package mocks

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
)

type MockCreateBookingUC struct {
	ExecuteFunc func(ctx context.Context, cmd *domain.BookingCommand) (*domain.BookingResponse, error)
}

func (m *MockCreateBookingUC) Execute(ctx context.Context, cmd *domain.BookingCommand) (*domain.BookingResponse, error) {
	return m.ExecuteFunc(ctx, cmd)
}
