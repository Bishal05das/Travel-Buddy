package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
)

type SearchHandler struct {
	uc port.Search
}

func NewSearchHandler(uc port.Search) *SearchHandler {
	return &SearchHandler{uc: uc}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	var filter domain.TourSearchFilter
	filter.Query = q.Get("q")

	// Malformed filters are rejected instead of silently becoming 0 / zero time.
	var err error
	if raw := q.Get("page"); raw != "" {
		filter.Page, err = strconv.Atoi(raw)
		if err != nil || filter.Page < 1 {
			http.Error(w, "page must be a positive integer", http.StatusBadRequest)
			return
		}
	}
	if raw := q.Get("limit"); raw != "" {
		filter.Limit, err = strconv.Atoi(raw)
		if err != nil || filter.Limit < 1 {
			http.Error(w, "limit must be a positive integer", http.StatusBadRequest)
			return
		}
	}
	if filter.MinPrice, err = parseOptionalFloat(q.Get("min_price")); err != nil {
		http.Error(w, "invalid min_price", http.StatusBadRequest)
		return
	}
	if filter.MaxPrice, err = parseOptionalFloat(q.Get("max_price")); err != nil {
		http.Error(w, "invalid max_price", http.StatusBadRequest)
		return
	}
	if filter.StartDate, err = parseOptionalDate(q.Get("start_date")); err != nil {
		http.Error(w, "invalid start_date, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	if filter.EndDate, err = parseOptionalDate(q.Get("end_date")); err != nil {
		http.Error(w, "invalid end_date, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	if err := filter.Prepare(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := h.uc.Execute(r.Context(), filter)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidSearchFilter) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	util.SendData(w, result, http.StatusOK)
}

func parseOptionalFloat(v string) (*float64, error) {
	if v == "" {
		return nil, nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func parseOptionalDate(v string) (*time.Time, error) {
	if v == "" {
		return nil, nil
	}
	t, err := time.Parse(time.DateOnly, v)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
