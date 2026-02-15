package domain

import "github.com/google/uuid"

type Customer struct {
	CustomerID uuid.UUID
	UserID     uuid.UUID
	Name       string
	Email      string
	Phone      int
}