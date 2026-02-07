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
	UpdateAgency(agency *domain.Agency) error
	DeleteAgency(agencyID int) error
}

type UserRepository interface {
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DeleteUser(userID int) error
}

type AgencyMemberRepository interface {
	CreateMember(member *domain.AgencyMember) error
	ListMember(agencyID int) ([]*domain.AgencyMember,error)
	UpdateMember(member *domain.AgencyMember) error
	DeleteMember(memberID int) error
}