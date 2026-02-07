package domain

import "time"

type Tour struct {
	ID                 int
	AgencyID           int
	Name               string
	StartDate          time.Time
	EndDate            time.Time
	AvailableSeat          int
	Description        string
	LastEnrollmentDate time.Time
	Price              int
	Discount           int
	Status             string
}
