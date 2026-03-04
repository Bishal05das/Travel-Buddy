package port

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/google/uuid"
)

type HomeRepository interface {
	GetTopTours(ctx context.Context, limit int) ([]domain.HomeTour, error)
	GetTopAgencies(ctx context.Context, limit int) ([]domain.HomeAgency, error)
}

type TourRepository interface {
	CreateTour(ctx context.Context, tour *domain.Tour) error
	ListTour(ctx context.Context, agencyID uuid.UUID) ([]*domain.Tour, error)
	UpdateTour(ctx context.Context, t *domain.Tour) error
	DeleteTour(ctx context.Context, tourID uuid.UUID) error
	GetByID(ctx context.Context, tourID uuid.UUID) (*domain.Tour, error)
	GetByIDForUpdate(ctx context.Context, tourID uuid.UUID) (*domain.Tour, error)
	UpdateAvailableSeats(ctx context.Context, tourID uuid.UUID, seats int) error
}

type AgencyRepository interface {
	CreateAgency(ctx context.Context, agency *domain.Agency) error
	UpdateAgency(ctx context.Context, agency *domain.Agency) error
	DeleteAgency(ctx context.Context, agencyID uuid.UUID) error
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	FindUser(ctx context.Context, email, pass string) (*domain.User, error)
}

type AgencyMemberRepository interface {
	CreateMember(ctx context.Context, member *domain.AgencyMember) error
	ListMember(ctx context.Context, agencyID uuid.UUID) ([]*domain.ListMemberResponse, error)
	UpdateMember(ctx context.Context, member *domain.AgencyMember) error
	DeleteMember(ctx context.Context, memberID uuid.UUID) error
	GetRoleIDFromMemberIDForUpdate(ctx context.Context,memberID uuid.UUID) (*int,error)
	FindMember(ctx context.Context,email,pass string) (*domain.AgencyMember,error)
}

type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BookingResponse, error)
	Update(ctx context.Context, booking *domain.Booking) error
	Cancel(ctx context.Context, id uuid.UUID) error
	GetOrCreateCustomerByUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	CreateCustomer(ctx context.Context, customer *domain.Customer) error
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
}

type RoleRepository interface {
	CreateRole(ctx context.Context, role *domain.Role) error
	DeleteRole(ctx context.Context, roleID int) error
	// CreateRolePermission(ctx context.Context, roleID int, permissionID int) error
	DeletePermissionsFromRole(ctx context.Context, roleID int) error
	AddPermissionsToRole(ctx context.Context, roleID int, permissionIDs []int) error
}

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}


type PermissionRepository interface {
	CreatePermission(ctx context.Context, permisson *domain.Permission) error
	DeletePermission(ctx context.Context,permissionID int) error
}