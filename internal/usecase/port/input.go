package port

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/google/uuid"
)

type CreateTour interface {
	Execute(tour *domain.Tour) error
}

type CreateBooking interface {
	Execute(ctx context.Context, req *domain.BookingRequest, userID *uuid.UUID, memberID *uuid.UUID) (*domain.BookingResponse, error)
}

type CreateUser interface {
	Execute(user *domain.User) error
}