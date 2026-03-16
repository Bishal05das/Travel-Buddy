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

func TestCreateAgencyHandler(t *testing.T) {

	tests := []struct {
		name           string
		body           string
		mockUsecase    func(*mocks.MockCreateAgency)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{
				"name":"TravelPro",
				"address":"Dhaka",
				"reg_id":"REG123"
			}`,
			mockUsecase: func(m *mocks.MockCreateAgency) {
				m.ExecuteFunc = func(ctx context.Context, a *domain.Agency) error {
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
		},

		{
			name:           "invalid json",
			body:           `{bad-json}`,
			mockUsecase:    func(m *mocks.MockCreateAgency) {},
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "validation error",
			body: `{
				"name":"",
				"address":"",
				"reg_id":""
			}`,
			mockUsecase:    func(m *mocks.MockCreateAgency) {},
			expectedStatus: http.StatusBadRequest,
		},

		{
			name: "usecase error",
			body: `{
				"name":"TravelPro",
				"address":"Dhaka",
				"reg_id":"REG123"
			}`,
			mockUsecase: func(m *mocks.MockCreateAgency) {
				m.ExecuteFunc = func(ctx context.Context, a *domain.Agency) error {
					return errors.New("database error")
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockCreateAgency{}
			tt.mockUsecase(mockUC)

			h := handler.NewAgencyHandler(mockUC, nil, nil)

			req := httptest.NewRequest(
				http.MethodPost,
				"/agency",
				bytes.NewBufferString(tt.body),
			)

			rec := httptest.NewRecorder()

			h.CreateAgency(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d",
					tt.expectedStatus,
					rec.Code)
			}
		})
	}
}

func TestUpdateAgencyHandler(t *testing.T) {

	tests := []struct {
		name           string
		agencyID       string
		body           string
		mockUsecase    func(*mocks.MockUpdateAgency)
		expectedStatus int
	}{
		{
			name:     "success",
			agencyID: "550e8400-e29b-41d4-a716-446655440000",
			body: `{
				"name":"Updated Travel",
				"address":"Dhaka",
				"reg_id":"REG123"
			}`,
			mockUsecase: func(m *mocks.MockUpdateAgency) {
				m.ExecuteFunc = func(ctx context.Context, a *domain.Agency) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
		},

		{
			name:           "invalid uuid",
			agencyID:       "bad-id",
			body:           `{}`,
			mockUsecase:    func(m *mocks.MockUpdateAgency) {},
			expectedStatus: http.StatusBadRequest,
		},

		{
			name:           "invalid json",
			agencyID:       "550e8400-e29b-41d4-a716-446655440000",
			body:           `{bad-json}`,
			mockUsecase:    func(m *mocks.MockUpdateAgency) {},
			expectedStatus: http.StatusBadRequest,
		},

		{
			name:     "validation error",
			agencyID: "550e8400-e29b-41d4-a716-446655440000",
			body: `{
				"name":"",
				"address":"",
				"reg_id":""
			}`,
			mockUsecase:    func(m *mocks.MockUpdateAgency) {},
			expectedStatus: http.StatusBadRequest,
		},

		{
			name:     "usecase error",
			agencyID: "550e8400-e29b-41d4-a716-446655440000",
			body: `{
				"name":"Travel",
				"address":"Dhaka",
				"reg_id":"REG123"
			}`,
			mockUsecase: func(m *mocks.MockUpdateAgency) {
				m.ExecuteFunc = func(ctx context.Context, a *domain.Agency) error {
					return errors.New("update failed")
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockUpdateAgency{}
			tt.mockUsecase(mockUC)

			h := handler.NewAgencyHandler(nil, mockUC, nil)

			req := httptest.NewRequest(
				http.MethodPut,
				"/agency/"+tt.agencyID,
				bytes.NewBufferString(tt.body),
			)

			req.SetPathValue("agency_id", tt.agencyID)

			rec := httptest.NewRecorder()

			h.UpdateAgency(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d",
					tt.expectedStatus,
					rec.Code)
			}
		})
	}
}

func TestDeleteAgencyHandler(t *testing.T) {

	tests := []struct {
		name           string
		agencyID       string
		mockUsecase    func(*mocks.MockDeleteAgency)
		expectedStatus int
	}{
		{
			name:     "success",
			agencyID: "550e8400-e29b-41d4-a716-446655440000",
			mockUsecase: func(m *mocks.MockDeleteAgency) {
				m.ExecuteFunc = func(ctx context.Context, id uuid.UUID) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
		},

		{
			name:           "invalid uuid",
			agencyID:       "bad-id",
			mockUsecase:    func(m *mocks.MockDeleteAgency) {},
			expectedStatus: http.StatusBadRequest,
		},

		{
			name:     "usecase error",
			agencyID: "550e8400-e29b-41d4-a716-446655440000",
			mockUsecase: func(m *mocks.MockDeleteAgency) {
				m.ExecuteFunc = func(ctx context.Context, id uuid.UUID) error {
					return errors.New("delete failed")
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockDeleteAgency{}
			tt.mockUsecase(mockUC)

			h := handler.NewAgencyHandler(nil, nil, mockUC)

			req := httptest.NewRequest(
				http.MethodDelete,
				"/agency/"+tt.agencyID,
				nil,
			)

			req.SetPathValue("agency_id", tt.agencyID)

			rec := httptest.NewRecorder()

			h.DeleteAgency(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d",
					tt.expectedStatus,
					rec.Code)
			}
		})
	}
}
