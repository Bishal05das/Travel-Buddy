package domain

import (
	"github.com/google/uuid"
)

type Payment struct {
	PaymentID     uuid.UUID
	BookingID     uuid.UUID
	TransactionID string
	Amount        float64
	Method        string
	Status        string
}
