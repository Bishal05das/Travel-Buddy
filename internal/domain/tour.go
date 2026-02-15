package domain

import (
	"time"

	"github.com/google/uuid"
)

type Tour struct {
	TourID             uuid.UUID
	AgencyID           int
	Name               string
	StartDate          time.Time
	EndDate            time.Time
	AvailableSeat      int
	Description        string
	LastEnrollmentDate time.Time
	Price              int
	Discount           int
	Status             string
}
