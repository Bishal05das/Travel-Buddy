package searchusecase_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	searchusecase "github.com/bishal05das/travelbuddy/internal/usecase/search"
)

type searchRepositoryStub struct {
	searchTours    func(context.Context, domain.TourSearchFilter) ([]domain.TourSearchResponse, error)
	countTours     func(context.Context, domain.TourSearchFilter) (int, error)
	searchAgencies func(context.Context, string, int) ([]domain.Agency, error)
}

func (r searchRepositoryStub) SearchTours(ctx context.Context, filter domain.TourSearchFilter) ([]domain.TourSearchResponse, error) {
	return r.searchTours(ctx, filter)
}

func (r searchRepositoryStub) CountTours(ctx context.Context, filter domain.TourSearchFilter) (int, error) {
	return r.countTours(ctx, filter)
}

func (r searchRepositoryStub) SearchAgencies(ctx context.Context, query string, limit int) ([]domain.Agency, error) {
	return r.searchAgencies(ctx, query, limit)
}

func TestSearchUseCaseNormalizesAndPaginates(t *testing.T) {
	tests := []struct {
		name      string
		filter    domain.TourSearchFilter
		count     int
		page      int
		limit     int
		offset    int
		totalPage int
	}{
		{"defaults and partial final page", domain.TourSearchFilter{Query: "  SAJEK\t Valley\n"}, 41, 1, 20, 0, 3},
		{"exact final page", domain.TourSearchFilter{Page: 2, Limit: 20}, 40, 2, 20, 20, 2},
		{"no matches", domain.TourSearchFilter{}, 0, 1, 20, 0, 0},
		{"requested page beyond results", domain.TourSearchFilter{Page: 5, Limit: 10}, 31, 5, 10, 40, 4},
		{"page size capped before offset", domain.TourSearchFilter{Page: 2, Limit: 100}, 101, 2, 50, 50, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			min, max := 1000.0, 4000.0
			start := time.Date(2031, 5, 10, 0, 0, 0, 0, time.UTC)
			end := start.AddDate(0, 0, 3)
			tt.filter.MinPrice, tt.filter.MaxPrice = &min, &max
			tt.filter.StartDate, tt.filter.EndDate = &start, &end
			wantFilter := tt.filter
			wantFilter.Query = strings.ToLower(strings.Join(strings.Fields(tt.filter.Query), " "))
			wantFilter.Page, wantFilter.Limit, wantFilter.Offset = tt.page, tt.limit, tt.offset
			wantTours := []domain.TourSearchResponse{{Name: "Sajek Valley", ImagePath: "images/tours/sajek.webp"}}
			wantAgencies := []domain.Agency{{Name: "Travel Agency"}}
			calls := []string{}
			repo := searchRepositoryStub{
				searchTours: func(gotCtx context.Context, gotFilter domain.TourSearchFilter) ([]domain.TourSearchResponse, error) {
					calls = append(calls, "tours")
					if gotCtx != ctx || !reflect.DeepEqual(gotFilter, wantFilter) {
						t.Errorf("SearchTours got context/filter %v %+v, want %+v", gotCtx, gotFilter, wantFilter)
					}
					return wantTours, nil
				},
				countTours: func(gotCtx context.Context, gotFilter domain.TourSearchFilter) (int, error) {
					calls = append(calls, "count")
					if gotCtx != ctx || !reflect.DeepEqual(gotFilter, wantFilter) {
						t.Errorf("CountTours got context/filter %v %+v, want %+v", gotCtx, gotFilter, wantFilter)
					}
					return tt.count, nil
				},
				searchAgencies: func(gotCtx context.Context, query string, limit int) ([]domain.Agency, error) {
					calls = append(calls, "agencies")
					if gotCtx != ctx || query != wantFilter.Query || limit != 5 {
						t.Errorf("SearchAgencies got context=%v query=%q limit=%d", gotCtx, query, limit)
					}
					return wantAgencies, nil
				},
			}
			got, err := searchusecase.NewSearchUseCase(repo).Execute(ctx, tt.filter)
			if err != nil {
				t.Fatal(err)
			}
			wantMeta := domain.SearchMeta{Page: tt.page, Limit: tt.limit, TotalCount: tt.count, TotalPage: tt.totalPage}
			if got.Meta != wantMeta || !reflect.DeepEqual(got.Tours, wantTours) || !reflect.DeepEqual(got.Agencies, wantAgencies) {
				t.Fatalf("unexpected result: %+v", got)
			}
			if !reflect.DeepEqual(calls, []string{"tours", "count", "agencies"}) {
				t.Fatalf("unexpected repository calls: %v", calls)
			}
		})
	}
}

func TestSearchUseCaseRejectsInvalidDirectFilters(t *testing.T) {
	ptr := func(value float64) *float64 { return &value }
	start := time.Date(2031, 5, 11, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, -1)
	status := "unknown"
	tests := []struct {
		name   string
		filter domain.TourSearchFilter
	}{
		{"NaN minimum", domain.TourSearchFilter{MinPrice: ptr(math.NaN())}},
		{"infinite maximum", domain.TourSearchFilter{MaxPrice: ptr(math.Inf(1))}},
		{"negative infinite minimum", domain.TourSearchFilter{MinPrice: ptr(math.Inf(-1))}},
		{"negative minimum", domain.TourSearchFilter{MinPrice: ptr(-1)}},
		{"negative maximum", domain.TourSearchFilter{MaxPrice: ptr(-1)}},
		{"inverted price range", domain.TourSearchFilter{MinPrice: ptr(101), MaxPrice: ptr(100)}},
		{"inverted date range", domain.TourSearchFilter{StartDate: &start, EndDate: &end}},
		{"negative page", domain.TourSearchFilter{Page: -1}},
		{"negative limit", domain.TourSearchFilter{Limit: -1}},
		{"page offset overflow", domain.TourSearchFilter{Page: math.MaxInt}},
		{"too many unicode characters", domain.TourSearchFilter{Query: strings.Repeat("ঢ", 201)}},
		{"too many terms", domain.TourSearchFilter{Query: "one two three four five six seven eight nine"}},
		{"invalid status", domain.TourSearchFilter{Status: &status}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := searchRepositoryStub{
				searchTours: func(context.Context, domain.TourSearchFilter) ([]domain.TourSearchResponse, error) {
					t.Fatal("invalid filter reached the repository")
					return nil, nil
				},
			}
			got, err := searchusecase.NewSearchUseCase(repo).Execute(context.Background(), tt.filter)
			if got != nil || !errors.Is(err, domain.ErrInvalidSearchFilter) {
				t.Fatalf("expected invalid filter error and nil result, got result=%+v error=%v", got, err)
			}
		})
	}
}

func TestSearchUseCasePropagatesRepositoryErrors(t *testing.T) {
	for _, stage := range []string{"tours", "count", "agencies"} {
		t.Run(stage, func(t *testing.T) {
			failure := errors.New("repository failure")
			calls := []string{}
			repo := searchRepositoryStub{
				searchTours: func(context.Context, domain.TourSearchFilter) ([]domain.TourSearchResponse, error) {
					calls = append(calls, "tours")
					if stage == "tours" {
						return nil, fmt.Errorf("search tours: %w", failure)
					}
					return []domain.TourSearchResponse{}, nil
				},
				countTours: func(context.Context, domain.TourSearchFilter) (int, error) {
					calls = append(calls, "count")
					if stage == "count" {
						return 0, fmt.Errorf("count tours: %w", failure)
					}
					return 0, nil
				},
				searchAgencies: func(context.Context, string, int) ([]domain.Agency, error) {
					calls = append(calls, "agencies")
					return nil, fmt.Errorf("search agencies: %w", failure)
				},
			}
			got, err := searchusecase.NewSearchUseCase(repo).Execute(context.Background(), domain.TourSearchFilter{})
			if got != nil || !errors.Is(err, failure) {
				t.Fatalf("expected repository error and nil result, got result=%+v error=%v", got, err)
			}
			wantCalls := map[string][]string{"tours": {"tours"}, "count": {"tours", "count"}, "agencies": {"tours", "count", "agencies"}}[stage]
			if !reflect.DeepEqual(calls, wantCalls) {
				t.Fatalf("repository continued after failure: calls=%v want=%v", calls, wantCalls)
			}
		})
	}
}
