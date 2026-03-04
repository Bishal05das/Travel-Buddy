package domain

import (
	"time"

	"github.com/google/uuid"
)

type HomeTour struct {
	TourID             uuid.UUID `json:"tour_id"`
	Name               string    `json:"name"`
	AgencyID           uuid.UUID `json:"agency_id"`
	AgencyName         string    `json:"agency_name"`
	AgencyRating       float64   `json:"agency_rating"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	Price              float64   `json:"price"`
	Discount           float64   `json:"discount"`
	FinalPrice         float64   `json:"final_price"`
	AvailableSeat      int       `json:"available_seat"`
	LastEnrollmentDate time.Time `json:"last_enrollment_date"`
	Description        string    `json:"description"`
	Status             string    `json:"status"`
	TotalBookings      int64     `json:"total_bookings"`
}

type HomeAgency struct {
	AgencyID   uuid.UUID `json:"agency_id"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	Rating     float64   `json:"rating"`
	TotalTours int64     `json:"total_tours"`
}

type HomeResponse struct {
	TopTours    []HomeTour   `json:"top_tours"`    // top 5 by bookings
	TopAgencies []HomeAgency `json:"top_agencies"` // top 6 by rating
}
