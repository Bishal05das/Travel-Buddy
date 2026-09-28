package port

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

type Home interface {
	GetHome(ctx context.Context) (*domain.HomeResponse, error)
}

type Search interface {
	Execute(ctx context.Context, filter domain.TourSearchFilter) (*domain.SearchResult, error)
}

type CreateTour interface {
	Execute(ctx context.Context, tour *domain.Tour) error
}

type ListTour interface {
	Execute(ctx context.Context, agencyID uuid.UUID, page, limit int) (*util.PaginationData, error)
}

type GetTour interface {
	Execute(ctx context.Context, tourID uuid.UUID) (*domain.Tour, error)
}

type UpdateTour interface {
	Execute(ctx context.Context, tour *domain.Tour) error
}

type DeleteTour interface {
	Execute(ctx context.Context, actor domain.Actor, tourID uuid.UUID) error
}

type UpdateTourStatus interface {
	// Cancelling a tour also cancels its active bookings.
	Execute(ctx context.Context, actor domain.Actor, tourID uuid.UUID, status string) (*domain.TourStatusChange, error)
}

type CreateBooking interface {
	Execute(ctx context.Context, req *domain.BookingCommand) (*domain.BookingResponse, error)
}

type ListBookings interface {
	Execute(ctx context.Context, actor domain.Actor, filter domain.BookingFilter) (*util.PaginationData, error)
}

type GetBooking interface {
	Execute(ctx context.Context, actor domain.Actor, scope domain.BookingScope, bookingID uuid.UUID) (*domain.BookingResponse, error)
}

// UpdateBookingStatus is the agency-side confirm / cancel / complete.
type UpdateBookingStatus interface {
	Execute(ctx context.Context, actor domain.Actor, agencyID, bookingID uuid.UUID, status string) (*domain.BookingResponse, error)
}

// CancelMyBooking lets a customer cancel their own booking.
type CancelMyBooking interface {
	Execute(ctx context.Context, actor domain.Actor, bookingID uuid.UUID) (*domain.BookingResponse, error)
}

type CreateUser interface {
	Execute(ctx context.Context, user *domain.User) error
}

type LoginUser interface {
	Execute(ctx context.Context, user *domain.ReqLogin) (*string, error)
}

type DeleteUser interface {
	Execute(ctx context.Context, userID uuid.UUID) error
}

type UpdateUser interface {
	Execute(ctx context.Context, user *domain.User) error
}

type CreateAgency interface {
	Execute(ctx context.Context, agency *domain.Agency, imagePath string) error
}

type UpdateAgency interface {
	Execute(ctx context.Context, agency *domain.Agency) error
}

type DeleteAgency interface {
	Execute(ctx context.Context, agencyID uuid.UUID) error
}

type CreateAgencyMember interface {
	Execute(ctx context.Context, actor domain.Actor, req *domain.CreateMemberRequest) error
}

type UpdateAgencyMemberPermission interface {
	Execute(ctx context.Context, actor domain.Actor, memberID uuid.UUID, req *domain.UpdatePermissionRequest) error
}

type DeleteAgencyMember interface {
	Execute(ctx context.Context, actor domain.Actor, agencyMemberID uuid.UUID) error
}

type ListAgencyMember interface {
	Execute(ctx context.Context, agencyID uuid.UUID) ([]*domain.ListMemberResponse, error)
}

type LoginMember interface {
	Execute(ctx context.Context, member *domain.ReqLogin) (*string, error)
}

type CreatePermission interface {
	Execute(ctx context.Context, p *domain.Permission) error
}

type DeletePermission interface {
	Execute(ctx context.Context, id int) error
}
