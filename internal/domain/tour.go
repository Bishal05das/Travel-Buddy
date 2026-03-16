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
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type CreateTourRequest struct {
	AgencyID           uuid.UUID `json:"agency_id" validate:"required,uuid"`
	Name               string    `json:"name" validate:"required,min=3,max=200"`
	StartDate          time.Time `json:"start_date" validate:"required"`
	EndDate            time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
	AvailableSeat      int       `json:"available_seat" validate:"required,gt=0"`
	Description        string    `json:"description" validate:"required,min=10,max=2000"`
	LastEnrollmentDate time.Time `json:"last_enrollment_date" validate:"required,ltefield=StartDate"`
	Price              float64   `json:"price" validate:"required,gt=0"`
	Discount           float64   `json:"discount" validate:"gte=0,lte=100"`
}

type UpdateTourRequest struct {
	AgencyID           uuid.UUID `json:"agency_id" validate:"required,uuid"`
	Name               string    `json:"name" validate:"required,min=3,max=200"`
	StartDate          time.Time `json:"start_date" validate:"required"`
	EndDate            time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
	AvailableSeat      int       `json:"available_seat" validate:"required,gt=0"`
	Description        string    `json:"description" validate:"required,min=10,max=2000"`
	LastEnrollmentDate time.Time `json:"last_enrollment_date" validate:"required,ltefield=StartDate"`
	Price              float64   `json:"price" validate:"required,gt=0"`
	Discount           float64   `json:"discount" validate:"gte=0,lte=100"`
}
