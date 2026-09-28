package domain

import (
	"errors"
	"fmt"
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

// BookingResponse is a booking as shown to the customer or the agency,
// including who booked, the tour and the payment.
type BookingResponse struct {
	BookingID      uuid.UUID  `json:"booking_id"`
	Status         string     `json:"status"`
	NumberOfPeople int        `json:"number_of_people"`
	TotalPrice     int        `json:"total_price"`
	BookingDate    time.Time  `json:"booking_date"`
	CreatedBy      string     `json:"created_by"` // "user" or "agency_member"
	UserID         *uuid.UUID `json:"user_id,omitempty"`
	MemberID       *uuid.UUID `json:"member_id,omitempty"`

	CustomerID    uuid.UUID `json:"customer_id"`
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	CustomerPhone string    `json:"customer_phone"`

	TourID        uuid.UUID `json:"tour_id"`
	TourName      string    `json:"tour_name"`
	TourStartDate time.Time `json:"tour_start_date"`
	AgencyID      uuid.UUID `json:"agency_id"`
	AgencyName    string    `json:"agency_name"`

	PaymentMethod string `json:"payment_method"`
	TransactionID string `json:"transaction_id"`
	PaymentStatus string `json:"payment_status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	BookingPending   = "pending"
	BookingConfirmed = "confirmed"
	BookingCancelled = "cancelled"
	BookingCompleted = "completed"
)

// bookingTransitions lists the statuses each status may move to.
// Cancelled and completed bookings are final.
var bookingTransitions = map[string][]string{
	BookingPending:   {BookingConfirmed, BookingCancelled},
	BookingConfirmed: {BookingCancelled, BookingCompleted},
}

var (
	ErrInvalidBookingTransition = errors.New("invalid booking status change")
	// ErrBookingNotFound also covers bookings outside the caller's scope, so
	// their existence is not revealed.
	ErrBookingNotFound = errors.New("booking not found")
)

// CheckBookingTransition reports whether a booking may move from -> to.
func CheckBookingTransition(from, to string) error {
	for _, allowed := range bookingTransitions[from] {
		if allowed == to {
			return nil
		}
	}
	return fmt.Errorf("%w: %s -> %s", ErrInvalidBookingTransition, from, to)
}

// BookingScope limits which bookings a query or change may touch: those of
// one agency's tours, or those made by one user. A nil field is unrestricted.
type BookingScope struct {
	AgencyID *uuid.UUID
	UserID   *uuid.UUID
}

// BookingFilter selects bookings for listing.
type BookingFilter struct {
	BookingScope
	TourID *uuid.UUID
	Status string
	Page   int
	Limit  int
}

// LockedBooking is a booking row locked for a status change, with the tour
// details the change needs.
type LockedBooking struct {
	BookingID      uuid.UUID
	TourID         uuid.UUID
	AgencyID       uuid.UUID
	UserID         *uuid.UUID
	NumberOfPeople int
	Status         string
	TourStartDate  time.Time
}

type UpdateBookingStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=confirmed cancelled completed"`
}
