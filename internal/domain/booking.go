package domain

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	BookingID      uuid.UUID
	CustomerID     uuid.UUID
	UserID         *uuid.UUID
	MemberID       *uuid.UUID
	TourID         uuid.UUID
	BookingDate    time.Time
	NumberOfPeople int
	TotalPrice     float64
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type BookingRequest struct {
	CustomerID     uuid.UUID `json:"customer_id" validate:"required,uuid"`
	TourID         uuid.UUID `json:"tour_id" validate:"required,uuid"`
	NumberOfPeople int       `json:"number_of_people" validate:"required,gt=0"`
	TotalPrice     int       `json:"total_price" validate:"required,gt=0"`
	Status         string    `json:"status" validate:"required,oneof=pending confirmed cancelled"`
	Method         string    `json:"method" validate:"required,oneof=cash card online"`
	TransactionId  string    `json:"transaction_id" validate:"omitempty,min=5,max=120"`

	CustomerName  string `json:"customer_name" validate:"required,min=2,max=120"`
	CustomerEmail string `json:"customer_email" validate:"required,email"`
	CustomerPhone string `json:"customer_phone" validate:"required,e164"`
}

type BookingResponse struct {
	BookingID      uuid.UUID
	CustomerID     uuid.UUID
	CustomerName   string
	TourID         uuid.UUID
	TourName       string
	AgencyName     string
	BookingDate    time.Time
	NumberOfPeople int
	TotalPrice     float64
	Status         string
	CreatedBy      string
	CreatedAt      time.Time
}
