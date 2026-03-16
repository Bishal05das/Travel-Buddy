package domain

import (
	"time"

	"github.com/google/uuid"
)

type Agency struct {
	AgencyID       uuid.UUID `json:"agency_id" db:"agency_id"`
	Name           string    `json:"name" db:"name"`
	Address        string    `json:"address" db:"address"`
	RegistrationID string    `json:"reg_id" db:"reg_id"`
	Rating         float64   `json:"rating"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

type CreateAgencyRequest struct {
	Name           string  `json:"name" validate:"required,min=3,max=150"`
	Address        string  `json:"address" validate:"required,min=5,max=300"`
	RegistrationID string  `json:"reg_id" validate:"required,alphanum,min=5,max=50"`
}

type UpdateAgencyRequest struct {
	Name           string  `json:"name" validate:"required,min=3,max=150"`
	Address        string  `json:"address" validate:"required,min=5,max=300"`
	RegistrationID string  `json:"reg_id" validate:"required,alphanum,min=5,max=50"`
}
