package mocks

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	util "github.com/bishal05das/travelbuddy/utils"
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
	ExecuteFunc func(ctx context.Context, agencyID uuid.UUID, page, limit int) (*util.PaginationData, error)
}

func (m *MockListTour) Execute(ctx context.Context, agencyID uuid.UUID, page, limit int) (*util.PaginationData, error) {
	return m.ExecuteFunc(ctx, agencyID, page, limit)
}

type MockUpdateTourStatus struct {
	ExecuteFunc func(ctx context.Context, actor domain.Actor, tourID uuid.UUID, status string) (*domain.TourStatusChange, error)
}

func (m *MockUpdateTourStatus) Execute(ctx context.Context, actor domain.Actor, tourID uuid.UUID, status string) (*domain.TourStatusChange, error) {
	return m.ExecuteFunc(ctx, actor, tourID, status)
}

type MockUpdateTour struct {
	ExecuteFunc func(ctx context.Context, tour *domain.Tour) error
}

func (m *MockUpdateTour) Execute(ctx context.Context, tour *domain.Tour) error {
	return m.ExecuteFunc(ctx, tour)
}

type MockDeleteTour struct {
	ExecuteFunc func(ctx context.Context, actor domain.Actor, id uuid.UUID) error
}

func (m *MockDeleteTour) Execute(ctx context.Context, actor domain.Actor, id uuid.UUID) error {
	return m.ExecuteFunc(ctx, actor, id)
}
