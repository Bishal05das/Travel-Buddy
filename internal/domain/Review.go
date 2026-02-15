package domain

import "github.com/google/uuid"

type Review struct {
	ReviewID uuid.UUID
	AgencyId uuid.UUID
	UserID   uuid.UUID
	Rating   int
	Comment  string
}
