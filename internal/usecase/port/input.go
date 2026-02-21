package port

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/google/uuid"
)

type CreateTour interface {
	Execute(ctx context.Context,tour *domain.Tour) error
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
	Execute(ctx context.Context, req *domain.UpdatePermissionRequest) error
}

type DeleteAgencyMember interface {
	Execute(ctx context.Context, agencyMemberID uuid.UUID) error
}

type ListAgencyMember interface {
	Execute(ctx context.Context, agencyID uuid.UUID) ([]*domain.AgencyMember, error)
}