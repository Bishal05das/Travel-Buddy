package handler_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
	"github.com/bishal05das/travelbuddy/internal/domain"
	mocks "github.com/bishal05das/travelbuddy/internal/mocks/usecase"
	"github.com/google/uuid"
)

func TestCreateMemberHandler(t *testing.T) {

	tests := []struct {
		name           string
		body           string
		mockUsecase    func(*mocks.MockCreateMember)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{
				"name":"john",
				"email":"john@test.com",
				"phone":"+8801700000000",
				"password":"12345678",
				"agency_id":"550e8400-e29b-41d4-a716-446655440000",
				"role_name":"admin",
				"permissions":[1,2,3]
			}`,
			mockUsecase: func(m *mocks.MockCreateMember) {
				m.ExecuteFunc = func(ctx context.Context, r *domain.CreateMemberRequest) error {
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json",
			body:           `{bad}`,
			mockUsecase:    func(m *mocks.MockCreateMember) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "validation error",
			body:           `{}`,
			mockUsecase:    func(m *mocks.MockCreateMember) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "usecase error",
			body: `{
				"name":"john",
				"email":"john@test.com",
				"password":"12345678",
				"agency_id":"550e8400-e29b-41d4-a716-446655440000"
			}`,
			mockUsecase: func(m *mocks.MockCreateMember) {
				m.ExecuteFunc = func(ctx context.Context, r *domain.CreateMemberRequest) error {
					return errors.New("failed")
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockCreateMember{}
			tt.mockUsecase(mockUC)

			h := handler.NewMemberHandler(mockUC, nil, nil, nil, nil)

			req := httptest.NewRequest(
				http.MethodPost,
				"/members",
				bytes.NewBufferString(tt.body),
			)

			rec := httptest.NewRecorder()

			h.CreateMember(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d",
					tt.expectedStatus,
					rec.Code)
			}
		})
	}
}

func TestDeleteMemberHandler(t *testing.T) {

	tests := []struct {
		name           string
		memberID       string
		mockUsecase    func(*mocks.MockDeleteMember)
		expectedStatus int
	}{
		{
			name:     "success",
			memberID: "550e8400-e29b-41d4-a716-446655440000",
			mockUsecase: func(m *mocks.MockDeleteMember) {
				m.ExecuteFunc = func(ctx context.Context, id uuid.UUID) error {
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid uuid",
			memberID:       "bad-id",
			mockUsecase:    func(m *mocks.MockDeleteMember) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "usecase error",
			memberID: "550e8400-e29b-41d4-a716-446655440000",
			mockUsecase: func(m *mocks.MockDeleteMember) {
				m.ExecuteFunc = func(ctx context.Context, id uuid.UUID) error {
					return errors.New("delete error")
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockDeleteMember{}
			tt.mockUsecase(mockUC)

			h := handler.NewMemberHandler(nil, mockUC, nil, nil, nil)

			req := httptest.NewRequest(http.MethodDelete, "/members/"+tt.memberID, nil)
			req.SetPathValue("member_id", tt.memberID)

			rec := httptest.NewRecorder()

			h.DeleteMember(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestListMemberHandler(t *testing.T) {

	tests := []struct {
		name           string
		agencyID       string
		mockUsecase    func(*mocks.MockListMember)
		expectedStatus int
	}{
		{
			name:     "success",
			agencyID: "550e8400-e29b-41d4-a716-446655440000",
			mockUsecase: func(m *mocks.MockListMember) {
				m.ExecuteFunc = func(ctx context.Context, agencyID uuid.UUID) ([]*domain.ListMemberResponse, error) {
					return []*domain.ListMemberResponse{
						{
							MemberID: agencyID,
							Name:     "member1",
							Email:    "member1@test.com",
							Phone:    "+8801700000000",
						},
						{
							MemberID: agencyID,
							Name:     "member2",
							Email:    "member2@test.com",
							Phone:    "+8801700000000",
						},
					}, nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid agency id",
			agencyID:       "bad-id",
			mockUsecase:    func(m *mocks.MockListMember) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "usecase error",
			agencyID: "550e8400-e29b-41d4-a716-446655440000",
			mockUsecase: func(m *mocks.MockListMember) {
				m.ExecuteFunc = func(ctx context.Context, agencyID uuid.UUID) ([]*domain.ListMemberResponse, error) {
					return nil, errors.New("failed to list members")
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockListMember{}
			tt.mockUsecase(mockUC)

			h := handler.NewMemberHandler(nil, nil, mockUC, nil, nil)

			req := httptest.NewRequest(
				http.MethodGet,
				"/agencies/"+tt.agencyID+"/members",
				nil,
			)

			req.SetPathValue("agency_id", tt.agencyID)

			rec := httptest.NewRecorder()

			h.ListMember(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestUpdateMemberPermissionsHandler(t *testing.T) {

	tests := []struct {
		name           string
		memberID       string
		body           string
		mockUsecase    func(*mocks.MockUpdatePermission)
		expectedStatus int
	}{
		{
			name:     "success",
			memberID: "550e8400-e29b-41d4-a716-446655440000",
			body: `{
				"permissions":[1,2,3]
			}`,
			mockUsecase: func(m *mocks.MockUpdatePermission) {
				m.ExecuteFunc = func(ctx context.Context, memberID uuid.UUID, req *domain.UpdatePermissionRequest) error {
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid member id",
			memberID:       "bad-id",
			body:           `{}`,
			mockUsecase:    func(m *mocks.MockUpdatePermission) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "invalid json",
			memberID: "550e8400-e29b-41d4-a716-446655440000",
			body:     `{bad-json}`,
			mockUsecase: func(m *mocks.MockUpdatePermission) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "usecase error",
			memberID: "550e8400-e29b-41d4-a716-446655440000",
			body: `{
				"permissions":["create_tour"]
			}`,
			mockUsecase: func(m *mocks.MockUpdatePermission) {
				m.ExecuteFunc = func(ctx context.Context, memberID uuid.UUID, req *domain.UpdatePermissionRequest) error {
					return errors.New("permission update failed")
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockUpdatePermission{}
			tt.mockUsecase(mockUC)

			h := handler.NewMemberHandler(nil, nil, nil, mockUC, nil)

			req := httptest.NewRequest(
				http.MethodPut,
				"/members/"+tt.memberID+"/permissions",
				bytes.NewBufferString(tt.body),
			)

			req.SetPathValue("member_id", tt.memberID)

			rec := httptest.NewRecorder()

			h.UpdateMemberPermissions(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestMemberLoginHandler(t *testing.T) {

	token := "jwt-token"

	tests := []struct {
		name           string
		body           string
		mockUsecase    func(*mocks.MockLoginMember)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{
				"email":"member@test.com",
				"password":"12345678"
			}`,
			mockUsecase: func(m *mocks.MockLoginMember) {
				m.ExecuteFunc = func(ctx context.Context, req *domain.ReqLogin) (*string, error) {
					return &token, nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json",
			body:           `{bad-json}`,
			mockUsecase:    func(m *mocks.MockLoginMember) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation error",
			body: `{}`,
			mockUsecase: func(m *mocks.MockLoginMember) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "usecase error",
			body: `{
				"email":"member@test.com",
				"password":"12345678"
			}`,
			mockUsecase: func(m *mocks.MockLoginMember) {
				m.ExecuteFunc = func(ctx context.Context, req *domain.ReqLogin) (*string, error) {
					return nil, errors.New("login failed")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockLoginMember{}
			tt.mockUsecase(mockUC)

			h := handler.NewMemberHandler(nil, nil, nil, nil, mockUC)

			req := httptest.NewRequest(
				http.MethodPost,
				"/members/login",
				bytes.NewBufferString(tt.body),
			)

			rec := httptest.NewRecorder()

			h.MemberLogin(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}