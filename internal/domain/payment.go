package domain

import (
	"github.com/google/uuid"
)

type Payment struct {
	PaymentID     int
	BookingID     uuid.UUID
	TransactionID string
	Amount        int
	Method        string
	Status        string
}
