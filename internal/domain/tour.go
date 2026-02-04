package domain

import "time"

type Tour struct {
	ID                 int
	AgencyID           int
	Name               string
	StartDate          time.Time
	EndDate            time.Time
	Description        string
	LastEnrollmentDate time.Time
	Price              int
	Discount           int
	Status             string
}

type TourRepository interface {
	CreateTour(tour *Tour) error
	ListTour(agencyID int) ([]*Tour,error)
	UpdateTour(t *Tour) error
	DeleteTour(tourID int) error
}
