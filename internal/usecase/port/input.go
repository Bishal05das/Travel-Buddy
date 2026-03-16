package port

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/google/uuid"
)

type Home interface {
	GetHome(ctx context.Context) (*domain.HomeResponse, error)
}

type Search interface {
	Execute(ctx context.Context, filter domain.TourSearchFilter) (*domain.SearchResult, error)
}

type CreateTour interface {
	Execute(ctx context.Context,tour *domain.Tour) error
}

type ListTour interface {
	Execute(ctx context.Context, agencyID uuid.UUID) ([]*domain.Tour, error)
}

type GetTour interface {
	Execute(ctx context.Context,tourID uuid.UUID) (*domain.Tour,error)
}

type UpdateTour interface {
	Execute(ctx context.Context, tour *domain.Tour) error
}

type DeleteTour interface {
	Execute(ctx context.Context, tourID uuid.UUID) error
}

type CreateBooking interface {
	Execute(ctx context.Context, req *domain.BookingRequest, userID *uuid.UUID, memberID *uuid.UUID) (*domain.BookingResponse, error)
}

type CreateUser interface {
	Execute(ctx context.Context,user *domain.User) error
}

type LoginUser interface {
	Execute(ctx context.Context,user *domain.ReqLogin) (*string, error)
}

type DeleteUser interface {
	Execute(ctx context.Context, userID uuid.UUID) error
}

type UpdateUser interface {
	Execute(ctx context.Context, user *domain.User) error
}

type CreateAgency interface {
	Execute(ctx context.Context, agency *domain.Agency) error
}

type UpdateAgency interface {
	Execute(ctx context.Context, agency *domain.Agency) error
}

type DeleteAgency interface {
	Execute(ctx context.Context, agencyID uuid.UUID) error
}

type CreateAgencyMember interface {
	Execute(ctx context.Context, req *domain.CreateMemberRequest) error
}

type UpdateAgencyMemberPermission interface {
	Execute(ctx context.Context,memberID uuid.UUID, req *domain.UpdatePermissionRequest) error
}

type DeleteAgencyMember interface {
	Execute(ctx context.Context, agencyMemberID uuid.UUID) error
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