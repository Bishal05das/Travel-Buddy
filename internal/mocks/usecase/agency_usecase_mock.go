package mocks

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/google/uuid"
)

type MockCreateAgency struct {
	ExecuteFunc func(ctx context.Context, agency *domain.Agency) error
}

func (m *MockCreateAgency) Execute(ctx context.Context, agency *domain.Agency) error {
	return m.ExecuteFunc(ctx, agency)
}

type MockUpdateAgency struct {
	ExecuteFunc func(ctx context.Context, agency *domain.Agency) error
}

func (m *MockUpdateAgency) Execute(ctx context.Context, agency *domain.Agency) error {
	return m.ExecuteFunc(ctx, agency)
}

type MockDeleteAgency struct {
	ExecuteFunc func(ctx context.Context, id uuid.UUID) error
}

func (m *MockDeleteAgency) Execute(ctx context.Context, id uuid.UUID) error {
	return m.ExecuteFunc(ctx, id)
}