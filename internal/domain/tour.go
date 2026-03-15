package domain

import (
	"time"

	"github.com/google/uuid"
)

type Tour struct {
	TourID             uuid.UUID `json:"tour_id" db:"tour_id"`
	AgencyID           uuid.UUID `json:"agency_id" db:"agency_id"`
	Name               string    `json:"name" db:"name"`
	StartDate          time.Time `json:"start_date" db:"start_date"`
	EndDate            time.Time `json:"end_date" db:"end_date"`
	AvailableSeat      int       `json:"available_seat" db:"available_seat"`
	Description        string    `json:"description" db:"description"`
	LastEnrollmentDate time.Time `json:"last_enrollment_date" db:"last_enrollment_date"`
	Price              float64   `json:"price" db:"price"`
	Discount           float64   `json:"discount" db:"discount"`
	Status             string    `json:"status" db:"status"`
}


