package domain

import "github.com/google/uuid"

type Agency struct {
	AgencyID       uuid.UUID
	Name           string
	Address        string
	RegistrationID string
	IsActive       bool
}
