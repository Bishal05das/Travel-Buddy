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

func TestCreateUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockUsecase    func(*mocks.MockCreateUser)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{
				"name":"bishal",
				"email":"test@test.com",
				"password":"12345678",
				"phone":"+8801700000000"
			}`,
			mockUsecase: func(m *mocks.MockCreateUser) {
				m.ExecuteFunc = func(ctx context.Context, u *domain.User) error {
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json",
			body:           `{invalid-json}`,
			mockUsecase:    func(m *mocks.MockCreateUser) {},
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "validation error",
			body: `{
				"name":"",
				"email":"invalid",
				"password":"1"
			}`,
			mockUsecase:    func(m *mocks.MockCreateUser) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "usecase failure",
			body: `{
				"name":"bishal",
				"email":"test@test.com",
				"password":"12345678",
				"phone":"+8801700000000"
			}`,
			mockUsecase: func(m *mocks.MockCreateUser) {
				m.ExecuteFunc = func(ctx context.Context, u *domain.User) error {
					return errors.New("db error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		mockUC := &mocks.MockCreateUser{}
		tt.mockUsecase(mockUC)

		h := handler.NewUserHandler(mockUC, nil, nil, nil)
		req := httptest.NewRequest(
			http.MethodPost,
			"/users",
			bytes.NewBufferString(tt.body),
		)
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()

		h.CreateUser(rec, req)

		if rec.Code != tt.expectedStatus {
			t.Errorf("expected status %d got %d",
				tt.expectedStatus,
				rec.Code,
			)
		}
	})
}
}

func TestUserLoginHandler(t *testing.T) {

	token := "jwt-token"

	tests := []struct {
		name           string
		body           string
		mockUsecase    func(*mocks.MockLoginUser)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{
				"email":"test@test.com",
				"password":"12345678"
			}`,
			mockUsecase: func(m *mocks.MockLoginUser) {
				m.ExecuteFunc = func(ctx context.Context, r *domain.ReqLogin) (*string, error) {
					return &token, nil
				}
			},
			expectedStatus: http.StatusCreated,
		},

		{
			name:           "invalid json",
			body:           `{bad}`,
			mockUsecase:    func(m *mocks.MockLoginUser) {},
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "usecase error",
			body: `{
				"email":"test@test.com",
				"password":"12345678"
			}`,
			mockUsecase: func(m *mocks.MockLoginUser) {
				m.ExecuteFunc = func(ctx context.Context, r *domain.ReqLogin) (*string, error) {
					return nil, errors.New("login failed")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockLoginUser{}
			tt.mockUsecase(mockUC)

			h := handler.NewUserHandler(nil, mockUC, nil, nil)

			req := httptest.NewRequest(
				http.MethodPost,
				"/login",
				bytes.NewBufferString(tt.body),
			)

			rec := httptest.NewRecorder()

			h.UserLogin(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestDeleteUserHandler(t *testing.T) {

	tests := []struct {
		name           string
		userID         string
		mockUsecase    func(*mocks.MockDeleteUser)
		expectedStatus int
	}{
		{
			name:   "success",
			userID: "550e8400-e29b-41d4-a716-446655440000",
			mockUsecase: func(m *mocks.MockDeleteUser) {
				m.ExecuteFunc = func(ctx context.Context, id uuid.UUID) error {
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
		},

		{
			name:           "invalid uuid",
			userID:         "bad-id",
			mockUsecase:    func(m *mocks.MockDeleteUser) {},
			expectedStatus: http.StatusBadRequest,
		},

		{
			name:   "usecase error",
			userID: "550e8400-e29b-41d4-a716-446655440000",
			mockUsecase: func(m *mocks.MockDeleteUser) {
				m.ExecuteFunc = func(ctx context.Context, id uuid.UUID) error {
					return errors.New("delete error")
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockDeleteUser{}
			tt.mockUsecase(mockUC)

			h := handler.NewUserHandler(nil, nil, mockUC, nil)

			req := httptest.NewRequest(
				http.MethodDelete,
				"/users/"+tt.userID,
				nil,
			)

			req.SetPathValue("user_id", tt.userID)

			rec := httptest.NewRecorder()

			h.DeleteUser(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
