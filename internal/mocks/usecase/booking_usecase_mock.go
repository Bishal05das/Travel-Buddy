package mocks

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/google/uuid"
)

type MockCreateBookingUC struct {
	ExecuteFunc func(
		ctx context.Context,
		req *domain.BookingRequest,
		userID *uuid.UUID,
		memberID *uuid.UUID,
	) (*domain.BookingResponse, error)
}

func (m *MockCreateBookingUC) Execute(
	ctx context.Context,
	req *domain.BookingRequest,
	userID *uuid.UUID,
	memberID *uuid.UUID,
) (*domain.BookingResponse, error) {
	return m.ExecuteFunc(ctx, req, userID, memberID)
}