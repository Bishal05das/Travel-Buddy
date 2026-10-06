package domain

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type TourSearchFilter struct {
	Query     string
	MinPrice  *float64
	MaxPrice  *float64
	StartDate *time.Time
	EndDate   *time.Time
	Status    *string
	AgencyID  *uuid.UUID
	Limit     int
	Offset    int
	Page      int
}

var ErrInvalidSearchFilter = errors.New("invalid search filter")

// NormalizeSearchQuery keeps direct API calls and both search forms consistent.
func NormalizeSearchQuery(query string) string {
	return strings.ToLower(strings.Join(strings.Fields(query), " "))
}

func (f *TourSearchFilter) Prepare() error {
	f.Query = NormalizeSearchQuery(f.Query)
	invalid := func(message string) error { return fmt.Errorf("%w: %s", ErrInvalidSearchFilter, message) }
	if utf8.RuneCountInString(f.Query) > 200 || len(strings.Fields(f.Query)) > 8 {
		return invalid("use at most 200 characters and 8 search words")
	}
	for _, price := range []*float64{f.MinPrice, f.MaxPrice} {
		if price != nil && (math.IsNaN(*price) || math.IsInf(*price, 0) || *price < 0) {
			return invalid("prices must be finite, non-negative numbers")
		}
	}
	if f.MinPrice != nil && f.MaxPrice != nil && *f.MinPrice > *f.MaxPrice {
		return invalid("minimum price cannot exceed maximum price")
	}
	if f.StartDate != nil && f.EndDate != nil && f.StartDate.After(*f.EndDate) {
		return invalid("start date cannot be after end date")
	}
	if f.Status != nil && *f.Status != TourOpen && *f.Status != TourClosed && *f.Status != TourCancelled {
		return invalid("invalid tour status")
	}
	if f.Page == 0 {
		f.Page = 1
	}
	if f.Limit == 0 {
		f.Limit = 20
	}
	if f.Page < 1 || f.Limit < 1 {
		return invalid("page and limit must be positive integers")
	}
	if f.Limit > 50 {
		f.Limit = 50
	}
	if f.Page-1 > math.MaxInt/f.Limit {
		return invalid("page is too large")
	}
	f.Offset = (f.Page - 1) * f.Limit
	return nil
}

type TourSearchResponse struct {
	TourID             uuid.UUID `json:"tour_id" db:"tour_id"`
	AgencyID           uuid.UUID `json:"agency_id" db:"agency_id"`
	AgencyName         string    `json:"agency_name"`
	Name               string    `json:"name" db:"name"`
	StartDate          time.Time `json:"start_date" db:"start_date"`
	EndDate            time.Time `json:"end_date" db:"end_date"`
	AvailableSeat      int       `json:"available_seat" db:"available_seat"`
	Description        string    `json:"description" db:"description"`
	LastEnrollmentDate time.Time `json:"last_enrollment_date" db:"last_enrollment_date"`
	Price              float64   `json:"price" db:"price"`
	Discount           float64   `json:"discount" db:"discount"`
	Status             string    `json:"status" db:"status"`
	ImagePath          string    `json:"image_path"`
}

type SearchResult struct {
	Tours    []TourSearchResponse
	Agencies []Agency
	Meta     SearchMeta
}

type SearchMeta struct {
	Page       int
	Limit      int
	TotalCount int
	TotalPage  int
}
