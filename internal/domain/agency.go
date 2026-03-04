package domain

import "github.com/google/uuid"

type Agency struct {
	AgencyID       uuid.UUID `json:"agency_id" db:"agency_id"`
	Name           string    `json:"name" db:"name"`
	Address        string    `json:"address" db:"address"`
	RegistrationID string    `json:"reg_id" db:"reg_id"`
	Rating         float64   `json:"rating"`
	IsActive       bool      `json:"is_active" db:"is_active"`
}
