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
	TotalPrice     int
	Status         string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type BookingCommand struct {
	TourID         uuid.UUID
	UserID         *uuid.UUID
	MemberID       *uuid.UUID
	AgencyID       *uuid.UUID // agency of the member creating a guest booking
	NumberOfPeople int
	TotalPrice     int
	Method         string
	TransactionId  string
	GuestInfo      *GuestCustomer
}

type GuestCustomer struct {
	Name  string
	Email string
	Phone string
}

type BookingRequestByUser struct {
	CustomerID     uuid.UUID `json:"customer_id" validate:"uuid"`
	TourID         uuid.UUID `json:"tour_id" validate:"required,uuid"`
	NumberOfPeople int       `json:"number_of_people" validate:"required,gt=0"`
	TotalPrice     int       `json:"total_price" validate:"required,gt=0"`
	Status         string    `json:"status" validate:"omitempty,oneof=pending confirmed cancelled"`
	Method         string    `json:"method" validate:"required,oneof=bkash bank nagad"`
	TransactionId  string    `json:"transaction_id" validate:"required,min=5,max=120"`
}

func (r *BookingRequestByUser) SetTourID(tourID uuid.UUID) {
	r.TourID = tourID
}

func (r *BookingRequestByUser) ToCommand(userID *uuid.UUID) *BookingCommand {
	return &BookingCommand{
		TourID:         r.TourID,
		UserID:         userID,
		NumberOfPeople: r.NumberOfPeople,
		Method:         r.Method,
		TransactionId:  r.TransactionId,
		TotalPrice:     r.TotalPrice,
	}
}

type BookingRequestByAdmin struct {
	CustomerID     uuid.UUID `json:"customer_id" validate:"uuid"`
	TourID         uuid.UUID `json:"tour_id" validate:"required,uuid"`
	NumberOfPeople int       `json:"number_of_people" validate:"required,gt=0"`
	TotalPrice     int       `json:"total_price" validate:"required,gt=0"`
	Status         string    `json:"status" validate:"omitempty,oneof=pending confirmed cancelled"`
	Method         string    `json:"method" validate:"required,oneof=bkash bank nagad"`
	TransactionId  string    `json:"transaction_id" validate:"required,min=5,max=120"`

	CustomerName  string `json:"customer_name" validate:"required,min=2,max=120"`
	CustomerEmail string `json:"customer_email" validate:"required,email"`
	CustomerPhone string `json:"customer_phone" validate:"required,e164"`
}

func (r *BookingRequestByAdmin) SetTourID(tourID uuid.UUID) {
	r.TourID = tourID
}

func (r *BookingRequestByAdmin) ToCommand(memberID *uuid.UUID) *BookingCommand {
	return &BookingCommand{
		TourID:         r.TourID,
		NumberOfPeople: r.NumberOfPeople,
		Method:         r.Method,
		TransactionId:  r.TransactionId,
		TotalPrice:     r.TotalPrice,
		MemberID:       memberID,
		GuestInfo: &GuestCustomer{
			Name:  r.CustomerName,
			Email: r.CustomerEmail,
			Phone: r.CustomerPhone,
		},
	}
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
