package mocks

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/google/uuid"
)

type MockCreateMember struct {
	ExecuteFunc func(ctx context.Context, actor domain.Actor, req *domain.CreateMemberRequest) error
}

func (m *MockCreateMember) Execute(ctx context.Context, actor domain.Actor, req *domain.CreateMemberRequest) error {
	return m.ExecuteFunc(ctx, actor, req)
}

type MockDeleteMember struct {
	ExecuteFunc func(ctx context.Context, actor domain.Actor, id uuid.UUID) error
}

func (m *MockDeleteMember) Execute(ctx context.Context, actor domain.Actor, id uuid.UUID) error {
	return m.ExecuteFunc(ctx, actor, id)
}

type MockListMember struct {
	ExecuteFunc func(ctx context.Context, agencyID uuid.UUID) ([]*domain.ListMemberResponse, error)
}

func (m *MockListMember) Execute(ctx context.Context, agencyID uuid.UUID) ([]*domain.ListMemberResponse, error) {
	return m.ExecuteFunc(ctx, agencyID)
}

type MockUpdatePermission struct {
	ExecuteFunc func(ctx context.Context, actor domain.Actor, memberID uuid.UUID, req *domain.UpdatePermissionRequest) error
}

func (m *MockUpdatePermission) Execute(ctx context.Context, actor domain.Actor, memberID uuid.UUID, req *domain.UpdatePermissionRequest) error {
	return m.ExecuteFunc(ctx, actor, memberID, req)
}

type MockLoginMember struct {
	ExecuteFunc func(ctx context.Context, req *domain.ReqLogin) (*string, error)
}

func (m *MockLoginMember) Execute(ctx context.Context, req *domain.ReqLogin) (*string, error) {
	return m.ExecuteFunc(ctx, req)
}
