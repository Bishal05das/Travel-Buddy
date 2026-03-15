package handler

import (
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

	query := r.URL.Query().Get("q")

	minPriceStr := r.URL.Query().Get("min_price")
	maxPriceStr := r.URL.Query().Get("max_price")

	var filter domain.TourSearchFilter
	filter.Query = query
	filter.Limit = 20
	filter.Offset = 0

	if minPriceStr != "" {
		v, _ := strconv.ParseFloat(minPriceStr, 64)
		filter.MinPrice = &v
	}

	if maxPriceStr != "" {
		v, _ := strconv.ParseFloat(maxPriceStr, 64)
		filter.MaxPrice = &v
	}

	startDateStr := r.URL.Query().Get("start_date")
	if startDateStr != "" {
		t, _ := time.Parse("2006-01-02", startDateStr)
		filter.StartDate = &t
	}

	result, err := h.uc.Execute(r.Context(), filter)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	util.SendData(w, result, http.StatusOK)
}
