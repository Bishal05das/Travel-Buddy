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

// type BookingRequest struct {
// 	TourID         uuid.UUID
// 	NumberOfPeople int
// 	CustomerName   string
// 	CustomerEmail  string
// 	CustomerPhone  int
// }
// type BookingUserRequest struct {
// 	CustomerID     uuid.UUID
// 	UserID         uuid.UUID
// 	TourID         uuid.UUID
// 	NumberOfPeople int
// 	TotalPrice     int
// 	Status         string
// }

type BookingRequest struct {
	CustomerID     uuid.UUID `json:"customer_id" db:"customer_id"`
	TourID         uuid.UUID `json:"tour_id" db:"tour_id"`
	NumberOfPeople int       `json:"number_of_people" db:"number_of_people"`
	TotalPrice     int       `json:"total_price" db:"total_price"`
	Status         string    `json:"status" db:"status"`
	Method         string    `json:"method" db:"method"`
	TransactionId  string    `json:"transaction_id" db:"transaction_id"`
	CustomerName   string    `json:"" db:""`
	CustomerEmail  string    `json:"" db:""`
	CustomerPhone  int       `json:"" db:""`
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
