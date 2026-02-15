package port

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/google/uuid"
)

type TourRepository interface {
	CreateTour(tour *domain.Tour) error
	ListTour(agencyID int) ([]*domain.Tour,error)
	UpdateTour(t *domain.Tour) error
	DeleteTour(tourID uuid.UUID) error
	GetByID(ctx context.Context,tourID uuid.UUID) (*domain.Tour,error)
	UpdateAvailableSeats(ctx context.Context, tourID uuid.UUID, seats int) error
}

type AgencyRepository interface {
	CreateAgency(agency *domain.Agency) error
	UpdateAgency(agency *domain.Agency) error
	DeleteAgency(agencyID int) error
}

type UserRepository interface {
	CreateUser(user *domain.User) error
	UpdateUser(user *domain.User) error
	DeleteUser(userID uuid.UUID) error
	FindUser(email,pass string) (*domain.User,error)
}

type AgencyMemberRepository interface {
	CreateMember(member *domain.AgencyMember) error
	ListMember(agencyID int) ([]*domain.AgencyMember,error)
	UpdateMember(member *domain.AgencyMember) error
	DeleteMember(memberID int) error
}

type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BookingResponse, error)
	Update(ctx context.Context, booking *domain.Booking) error
	Cancel(ctx context.Context, id uuid.UUID) error
	GetOrCreateCustomerByUser(ctx context.Context, userID uuid.UUID) (uuid.UUID,error)
	CreateCustomer(ctx context.Context, customer *domain.Customer) error
}

type PaymentRepository interface {
	Create(payment *domain.Payment) error
}