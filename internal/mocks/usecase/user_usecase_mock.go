package mocks

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/google/uuid"
)

type MockCreateUser struct {
	ExecuteFunc func(ctx context.Context, user *domain.User) error
}

func (m *MockCreateUser) Execute(ctx context.Context, user *domain.User) error {
	return m.ExecuteFunc(ctx, user)
}

type MockLoginUser struct {
	ExecuteFunc func(ctx context.Context, req *domain.ReqLogin) (*string, error)
}

func (m *MockLoginUser) Execute(ctx context.Context, req *domain.ReqLogin) (*string, error) {
	return m.ExecuteFunc(ctx, req)
}

type MockDeleteUser struct {
	ExecuteFunc func(ctx context.Context, id uuid.UUID) error
}

func (m *MockDeleteUser) Execute(ctx context.Context, id uuid.UUID) error {
	return m.ExecuteFunc(ctx, id)
}

type MockUpdateUser struct {
	ExecuteFunc func(ctx context.Context, user *domain.User) error
}

func (m *MockUpdateUser) Execute(ctx context.Context, user *domain.User) error {
	return m.ExecuteFunc(ctx, user)
}
