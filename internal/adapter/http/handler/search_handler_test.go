package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
	"github.com/bishal05das/travelbuddy/internal/domain"
)

type searchUseCaseStub struct {
	execute func(context.Context, domain.TourSearchFilter) (*domain.SearchResult, error)
}

func (s searchUseCaseStub) Execute(ctx context.Context, filter domain.TourSearchFilter) (*domain.SearchResult, error) {
	return s.execute(ctx, filter)
}

func TestSearchHandlerPreparesFilters(t *testing.T) {
	tests := []struct {
		name   string
		params url.Values
		query  string
		page   int
		limit  int
		offset int
	}{
		{name: "default pagination", params: url.Values{}, page: 1, limit: 20},
		{name: "normalized keyword", params: url.Values{"q": {"  SAJEK\t  Valley\n "}}, query: "sajek valley", page: 1, limit: 20},
		{name: "whitespace only", params: url.Values{"q": {" \t\n\u00a0"}}, page: 1, limit: 20},
		{name: "custom pagination", params: url.Values{"page": {"3"}, "limit": {"7"}}, page: 3, limit: 7, offset: 14},
		{name: "bounded page size", params: url.Values{"page": {"3"}, "limit": {"100"}}, page: 3, limit: 50, offset: 100},
		{name: "200 unicode characters", params: url.Values{"q": {strings.Repeat("ঢ", 200)}}, query: strings.Repeat("ঢ", 200), page: 1, limit: 20},
		{name: "eight words", params: url.Values{"q": {"one two three four five six seven eight"}}, query: "one two three four five six seven eight", page: 1, limit: 20},
		{name: "inclusive zero price and equal dates", params: url.Values{"min_price": {"0"}, "max_price": {"0"}, "start_date": {"2031-05-10"}, "end_date": {"2031-05-10"}}, page: 1, limit: 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			uc := searchUseCaseStub{execute: func(_ context.Context, filter domain.TourSearchFilter) (*domain.SearchResult, error) {
				calls++
				if filter.Query != tt.query || filter.Page != tt.page || filter.Limit != tt.limit || filter.Offset != tt.offset {
					t.Errorf("unexpected prepared filter: %+v", filter)
				}
				if tt.params.Has("min_price") && (filter.MinPrice == nil || *filter.MinPrice != 0 || filter.MaxPrice == nil || *filter.MaxPrice != 0) {
					t.Errorf("zero price was not preserved: %+v", filter)
				}
				if tt.params.Has("start_date") && (filter.StartDate == nil || filter.EndDate == nil || !filter.StartDate.Equal(*filter.EndDate) || filter.StartDate.Format("2006-01-02") != "2031-05-10") {
					t.Errorf("date bounds were not preserved: %+v", filter)
				}
				return &domain.SearchResult{Tours: []domain.TourSearchResponse{}, Agencies: []domain.Agency{}}, nil
			}}
			rec := httptest.NewRecorder()
			handler.NewSearchHandler(uc).Search(rec, httptest.NewRequest(http.MethodGet, "/search?"+tt.params.Encode(), nil))
			if rec.Code != http.StatusOK || calls != 1 {
				t.Fatalf("status=%d calls=%d body=%q", rec.Code, calls, rec.Body.String())
			}
		})
	}
}

func TestSearchHandlerRejectsInvalidFilters(t *testing.T) {
	tests := []struct {
		name   string
		params url.Values
	}{
		{"malformed minimum", url.Values{"min_price": {"abc"}}},
		{"malformed maximum", url.Values{"max_price": {"abc"}}},
		{"negative minimum", url.Values{"min_price": {"-1"}}},
		{"negative maximum", url.Values{"max_price": {"-1"}}},
		{"NaN minimum", url.Values{"min_price": {"NaN"}}},
		{"NaN maximum", url.Values{"max_price": {"NaN"}}},
		{"infinite minimum", url.Values{"min_price": {"+Inf"}}},
		{"infinite maximum", url.Values{"max_price": {"Inf"}}},
		{"negative infinity", url.Values{"max_price": {"-Inf"}}},
		{"price overflow", url.Values{"max_price": {"1e999"}}},
		{"inverted prices", url.Values{"min_price": {"101"}, "max_price": {"100"}}},
		{"invalid start date", url.Values{"start_date": {"2031-02-29"}}},
		{"invalid end date", url.Values{"end_date": {"not-a-date"}}},
		{"timestamp rather than date", url.Values{"start_date": {"2031-05-10T00:00:00Z"}}},
		{"inverted dates", url.Values{"start_date": {"2031-05-11"}, "end_date": {"2031-05-10"}}},
		{"malformed page", url.Values{"page": {"many"}}},
		{"fractional page", url.Values{"page": {"1.5"}}},
		{"zero page", url.Values{"page": {"0"}}},
		{"negative page", url.Values{"page": {"-1"}}},
		{"page parsing overflow", url.Values{"page": {"999999999999999999999999999999999"}}},
		{"page offset overflow", url.Values{"page": {strconv.Itoa(math.MaxInt)}}},
		{"malformed limit", url.Values{"limit": {"many"}}},
		{"fractional limit", url.Values{"limit": {"1.5"}}},
		{"zero limit", url.Values{"limit": {"0"}}},
		{"negative limit", url.Values{"limit": {"-1"}}},
		{"too many unicode characters", url.Values{"q": {strings.Repeat("ঢ", 201)}}},
		{"too many terms", url.Values{"q": {"one two three four five six seven eight nine"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := searchUseCaseStub{execute: func(context.Context, domain.TourSearchFilter) (*domain.SearchResult, error) {
				t.Fatal("invalid request reached the use case")
				return nil, nil
			}}
			rec := httptest.NewRecorder()
			handler.NewSearchHandler(uc).Search(rec, httptest.NewRequest(http.MethodGet, "/search?"+tt.params.Encode(), nil))
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected HTTP 400, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestSearchHandlerUseCaseErrors(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
		body   string
	}{
		{"invalid filter", fmt.Errorf("%w: invalid range", domain.ErrInvalidSearchFilter), http.StatusBadRequest, "invalid range"},
		{"repository failure", errors.New("database secret must not be returned"), http.StatusInternalServerError, "internal server error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := searchUseCaseStub{execute: func(context.Context, domain.TourSearchFilter) (*domain.SearchResult, error) { return nil, tt.err }}
			rec := httptest.NewRecorder()
			handler.NewSearchHandler(uc).Search(rec, httptest.NewRequest(http.MethodGet, "/search", nil))
			if rec.Code != tt.status || !strings.Contains(rec.Body.String(), tt.body) {
				t.Fatalf("unexpected response: status=%d body=%q", rec.Code, rec.Body.String())
			}
			if strings.Contains(rec.Body.String(), "database secret") {
				t.Fatal("response exposes an internal error")
			}
		})
	}
}

func TestSearchHandlerReturnsPaginationAndImages(t *testing.T) {
	want := domain.SearchResult{
		Tours:    []domain.TourSearchResponse{{Name: "Sajek Valley", ImagePath: "images/tours/sajek.webp"}},
		Agencies: []domain.Agency{},
		Meta:     domain.SearchMeta{Page: 2, Limit: 20, TotalCount: 41, TotalPage: 3},
	}
	uc := searchUseCaseStub{execute: func(context.Context, domain.TourSearchFilter) (*domain.SearchResult, error) { return &want, nil }}
	rec := httptest.NewRecorder()
	handler.NewSearchHandler(uc).Search(rec, httptest.NewRequest(http.MethodGet, "/search?page=2", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status %d: %s", rec.Code, rec.Body.String())
	}
	var got domain.SearchResult
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if got.Meta != want.Meta || len(got.Tours) != 1 || got.Tours[0].ImagePath != want.Tours[0].ImagePath {
		t.Fatalf("response omitted pagination or image: %+v", got)
	}
}
