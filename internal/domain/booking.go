package domain

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	BookingID      uuid.UUID
	CustomerID     uuid.UUID
	UserID         uuid.UUID
	MemberID       uuid.UUID
	TourID         uuid.UUID
	BookingDate    time.Time
	NumberOfPeople int
	TotalPrice     int
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
	CustomerID     uuid.UUID
	TourID         uuid.UUID
	NumberOfPeople int
	TotalPrice     int
	Status         string
	CustomerName   string
	CustomerEmail  string
	CustomerPhone  int
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
