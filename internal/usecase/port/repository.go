package port

import "github.com/bishal05das/travelbuddy/internal/domain"

type TourRepository interface {
	CreateTour(tour *domain.Tour) error
	ListTour(agencyID int) ([]*domain.Tour,error)
	UpdateTour(t *domain.Tour) error
	DeleteTour(tourID int) error
}

type AgencyRepository interface {
	CreateAgency(agency *domain.Agency) error
	UpdateAgency(t *domain.Agency) error
	DeleteAgency(agencyID int) error
}

type UserRepository interface {
	CreateUser(user *domain.User) error
	UpdateUser(t *domain.User) error
	DeleteUser(userID int) error
}