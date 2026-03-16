package mocks

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/google/uuid"
)

type MockCreateTour struct {
	ExecuteFunc func(ctx context.Context, tour *domain.Tour) error
}

func (m *MockCreateTour) Execute(ctx context.Context, tour *domain.Tour) error {
	return m.ExecuteFunc(ctx, tour)
}

type MockGetTour struct {
	ExecuteFunc func(ctx context.Context, id uuid.UUID) (*domain.Tour, error)
}

func (m *MockGetTour) Execute(ctx context.Context, id uuid.UUID) (*domain.Tour, error) {
	return m.ExecuteFunc(ctx, id)
}

type MockListTour struct {
	ExecuteFunc func(ctx context.Context, agencyID uuid.UUID) ([]*domain.Tour, error)
}

func (m *MockListTour) Execute(ctx context.Context, agencyID uuid.UUID) ([]*domain.Tour, error) {
	return m.ExecuteFunc(ctx, agencyID)
}

type MockUpdateTour struct {
	ExecuteFunc func(ctx context.Context, tour *domain.Tour) error
}

func (m *MockUpdateTour) Execute(ctx context.Context, tour *domain.Tour) error {
	return m.ExecuteFunc(ctx, tour)
}

type MockDeleteTour struct {
	ExecuteFunc func(ctx context.Context, id uuid.UUID) error
}

func (m *MockDeleteTour) Execute(ctx context.Context, id uuid.UUID) error {
	return m.ExecuteFunc(ctx, id)
}