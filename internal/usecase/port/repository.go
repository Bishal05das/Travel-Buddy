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
	ListTour(ctx context.Context, agencyID uuid.UUID, page, limit int) ([]*domain.Tour, error)
	Count(ctx context.Context, agencyID uuid.UUID) (int, error)
	UpdateTour(ctx context.Context, t *domain.Tour) error
	// agencyScope, when non-nil, restricts the change to tours of that agency.
	DeleteTour(ctx context.Context, tourID uuid.UUID, agencyScope *uuid.UUID) error
	GetByID(ctx context.Context, tourID uuid.UUID) (*domain.Tour, error)
	GetByIDForUpdate(ctx context.Context, tourID uuid.UUID) (*domain.Tour, error)
	UpdateAvailableSeats(ctx context.Context, tourID uuid.UUID, seats int) error
	UpdateTourStatus(ctx context.Context, tourID uuid.UUID, status string, agencyScope *uuid.UUID) error
}

type AgencyRepository interface {
	CreateAgency(ctx context.Context, agency *domain.Agency, imagePath string) error
	UpdateAgency(ctx context.Context, agency *domain.Agency) error
	DeleteAgency(ctx context.Context, agencyID uuid.UUID) error
	UpdateAgencyImage(ctx context.Context, agencyID uuid.UUID, newImagePath string) (oldImagePath string, err error)
	GetCurrentImage(ctx context.Context, agencyID uuid.UUID) (*domain.AgencyImage, error)
	ListAgencyImages(ctx context.Context, agencyID uuid.UUID) ([]*domain.AgencyImage, error)
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
}

type AgencyMemberRepository interface {
	CreateMember(ctx context.Context, member *domain.AgencyMember) error
	ListMember(ctx context.Context, agencyID uuid.UUID) ([]*domain.ListMemberResponse, error)
	UpdateMember(ctx context.Context, member *domain.AgencyMember) error
	DeleteMember(ctx context.Context, memberID uuid.UUID, agencyScope *uuid.UUID) error
	GetRoleIDFromMemberIDForUpdate(ctx context.Context, memberID uuid.UUID, agencyScope *uuid.UUID) (*int, error)
	GetPermissionIDs(ctx context.Context, memberID uuid.UUID) ([]int, error)
	FindMember(ctx context.Context, email string) (*domain.AgencyMember, error)
}

type BookingRepository interface {
	Create(ctx context.Context, booking *domain.Booking) error
	GetByID(ctx context.Context, id uuid.UUID, scope domain.BookingScope) (*domain.BookingResponse, error)
	// List returns one page of bookings and the total number matching.
	List(ctx context.Context, filter domain.BookingFilter) ([]*domain.BookingResponse, int, error)
	// GetForUpdate locks the booking row (within scope) for a status change.
	GetForUpdate(ctx context.Context, id uuid.UUID, scope domain.BookingScope) (*domain.LockedBooking, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	GetOrCreateCustomerByUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error)
	CreateCustomer(ctx context.Context, customer *domain.Customer) error
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *domain.Payment) error
	// SetStatusForBooking moves the booking's payments from one status to another.
	SetStatusForBooking(ctx context.Context, bookingID uuid.UUID, from, to string) error
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
	DeletePermission(ctx context.Context, permissionID int) error
}

// Authorizer answers permission checks for agency members.
type Authorizer interface {
	MemberHasPermission(ctx context.Context, memberID, agencyID uuid.UUID, resource, action string) (bool, error)
}

type SearchRepository interface {
	SearchTours(ctx context.Context, filter domain.TourSearchFilter) ([]domain.TourSearchResponse, error)
	SearchAgencies(ctx context.Context, query string, limit int) ([]domain.Agency, error)
}
