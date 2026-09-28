package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Tour struct {
	TourID             uuid.UUID `json:"tour_id" db:"tour_id"`
	AgencyID           uuid.UUID `json:"agency_id" db:"agency_id"`
	Name               string    `json:"name" db:"name"`
	StartDate          time.Time `json:"start_date" db:"start_date"`
	EndDate            time.Time `json:"end_date" db:"end_date"`
	TotalSeat          int       `json:"total_seat" db:"total_seat"`         // capacity, set by the agency
	AvailableSeat      int       `json:"available_seat" db:"available_seat"` // capacity minus active bookings
	Description        string    `json:"description" db:"description"`
	LastEnrollmentDate time.Time `json:"last_enrollment_date" db:"last_enrollment_date"`
	Price              int       `json:"price" db:"price"`
	Discount           int       `json:"discount" db:"discount"`
	Status             string    `json:"status" db:"status"`
	ImagePath          string    `json:"image_path" db:"image_path"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// Validate enforces the tour invariants independently of the HTTP layer.
func (t *Tour) Validate() error {
	switch {
	case strings.TrimSpace(t.Name) == "":
		return errors.New("tour name is required")
	case !t.EndDate.After(t.StartDate):
		return errors.New("end date must be after start date")
	case t.LastEnrollmentDate.After(t.StartDate):
		return errors.New("last enrollment date must not be after start date")
	case t.TotalSeat < 1:
		return errors.New("total seats must be at least 1")
	case t.Price <= 0:
		return errors.New("price must be greater than 0")
	case t.Discount < 0 || t.Discount > 100:
		return errors.New("discount must be between 0 and 100")
	}
	return nil
}

// BookedSeats is the number of seats held by active bookings.
func (t *Tour) BookedSeats() int {
	return t.TotalSeat - t.AvailableSeat
}

// ErrCapacityBelowBooked is returned when an update would shrink a tour
// below the seats that are already booked.
var ErrCapacityBelowBooked = errors.New("total seats cannot be lower than the seats already booked")

// UnitPrice is the per-person price after applying the percentage discount.
func (t *Tour) UnitPrice() int {
	return t.Price - (t.Price*t.Discount)/100
}

type CreateTourRequest struct {
	AgencyID           uuid.UUID `json:"agency_id" validate:"required,uuid"`
	Name               string    `json:"name" validate:"required,min=3,max=200"`
	StartDate          time.Time `json:"start_date" validate:"required"`
	EndDate            time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
	TotalSeat          int       `json:"total_seat" validate:"required,gt=0"`
	Description        string    `json:"description" validate:"required,min=10,max=2000"`
	LastEnrollmentDate time.Time `json:"last_enrollment_date" validate:"required,ltefield=StartDate"`
	Price              int       `json:"price" validate:"required,gt=0"`
	Discount           int       `json:"discount" validate:"gte=0,lte=100"`
}

type UpdateTourRequest struct {
	AgencyID           uuid.UUID `json:"agency_id" validate:"required,uuid"`
	Name               string    `json:"name" validate:"required,min=3,max=200"`
	StartDate          time.Time `json:"start_date" validate:"required"`
	EndDate            time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
	TotalSeat          int       `json:"total_seat" validate:"required,gt=0"`
	Description        string    `json:"description" validate:"required,min=10,max=2000"`
	LastEnrollmentDate time.Time `json:"last_enrollment_date" validate:"required,ltefield=StartDate"`
	Price              int       `json:"price" validate:"required,gt=0"`
	Discount           int       `json:"discount" validate:"gte=0,lte=100"`
}
