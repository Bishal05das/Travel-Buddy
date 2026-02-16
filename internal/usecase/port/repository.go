package port

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/google/uuid"
)

type TourRepository interface {
	CreateTour(ctx context.Context,tour *domain.Tour) error
	ListTour(ctx context.Context,agencyID int) ([]*domain.Tour,error)
	UpdateTour(ctx context.Context,t *domain.Tour) error
	DeleteTour(ctx context.Context,tourID uuid.UUID) error
	GetByID(ctx context.Context,tourID uuid.UUID) (*domain.Tour,error)
	UpdateAvailableSeats(ctx context.Context, tourID uuid.UUID, seats int) error
}

type AgencyRepository interface {
	CreateAgency(ctx context.Context,agency *domain.Agency) error
	UpdateAgency(ctx context.Context,agency *domain.Agency) error
	DeleteAgency(ctx context.Context,agencyID int) error
}

type UserRepository interface {
	CreateUser(ctx context.Context,user *domain.User) error
	UpdateUser(ctx context.Context,user *domain.User) error
	DeleteUser(ctx context.Context,userID uuid.UUID) error
	FindUser(ctx context.Context,email,pass string) (*domain.User,error)
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
	Create(ctx context.Context,payment *domain.Payment) error
}