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

func TestCreateTourHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockUsecase    func(*mocks.MockCreateTour)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{
				"agency_id":"550e8400-e29b-41d4-a716-446655440000",
				"name":"cox tour",
				"start_date":"2026-06-01T00:00:00Z",
				"end_date":"2026-06-10T00:00:00Z",
				"last_enrollment_date":"2026-05-25T00:00:00Z",
				"available_seat":10,
				"description":"this is a test tour description",
				"price":1000,
				"discount":0
			}`,
			mockUsecase: func(m *mocks.MockCreateTour) {
				m.ExecuteFunc = func(ctx context.Context, tour *domain.Tour) error {
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid json",
			body:           `{bad-json}`,
			mockUsecase:    func(m *mocks.MockCreateTour) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "validation error",
			body:           `{}`,
			mockUsecase:    func(m *mocks.MockCreateTour) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "usecase error",
			body: `{
				"agency_id":"550e8400-e29b-41d4-a716-446655440000",
				"name":"cox tour",
				"start_date":"2026-06-01T00:00:00Z",
				"end_date":"2026-06-10T00:00:00Z",
				"last_enrollment_date":"2026-05-25T00:00:00Z",
				"available_seat":10,
				"description":"this is a test tour description",
				"price":1000,
				"discount":0
			}`,
			mockUsecase: func(m *mocks.MockCreateTour) {
				m.ExecuteFunc = func(ctx context.Context, tour *domain.Tour) error {
					return errors.New("db error")
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUC := &mocks.MockCreateTour{}
			tt.mockUsecase(mockUC)
			h := handler.NewTourHandler(mockUC, nil, nil, nil, nil)
			req := httptest.NewRequest(
				http.MethodPost,
				"/tours",
				bytes.NewBufferString(tt.body),
			)
			rec := httptest.NewRecorder()
			h.Create(rec, req)
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestGetTourHandler(t *testing.T) {

	tests := []struct {
		name           string
		tourID         string
		mockUsecase    func(*mocks.MockGetTour)
		expectedStatus int
	}{
		{
			name:   "success",
			tourID: "550e8400-e29b-41d4-a716-446655440000",
			mockUsecase: func(m *mocks.MockGetTour) {
				m.ExecuteFunc = func(ctx context.Context, id uuid.UUID) (*domain.Tour, error) {
					return &domain.Tour{Name: "cox tour"}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid uuid",
			tourID:         "bad-id",
			mockUsecase:    func(m *mocks.MockGetTour) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "usecase error",
			tourID: "550e8400-e29b-41d4-a716-446655440000",
			mockUsecase: func(m *mocks.MockGetTour) {
				m.ExecuteFunc = func(ctx context.Context, id uuid.UUID) (*domain.Tour, error) {
					return nil, errors.New("not found")
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockGetTour{}
			tt.mockUsecase(mockUC)

			h := handler.NewTourHandler(nil, mockUC, nil, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/tours/"+tt.tourID, nil)
			req.SetPathValue("tour_id", tt.tourID)

			rec := httptest.NewRecorder()

			h.Get(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestListTourHandler(t *testing.T) {

	tests := []struct {
		name           string
		agencyID       string
		mockUsecase    func(*mocks.MockListTour)
		expectedStatus int
	}{
		{
			name:     "success",
			agencyID: "550e8400-e29b-41d4-a716-446655440000",
			mockUsecase: func(m *mocks.MockListTour) {
				m.ExecuteFunc = func(ctx context.Context, id uuid.UUID) ([]*domain.Tour, error) {
					return []*domain.Tour{{Name: "cox tour"}}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid agency id",
			agencyID:       "bad-id",
			mockUsecase:    func(m *mocks.MockListTour) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "usecase error",
			agencyID: "550e8400-e29b-41d4-a716-446655440000",
			mockUsecase: func(m *mocks.MockListTour) {
				m.ExecuteFunc = func(ctx context.Context, id uuid.UUID) ([]*domain.Tour, error) {
					return nil, errors.New("db error")
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockListTour{}
			tt.mockUsecase(mockUC)

			h := handler.NewTourHandler(nil, nil, mockUC, nil, nil)

			req := httptest.NewRequest(http.MethodGet, "/agencies/"+tt.agencyID+"/tours", nil)
			req.SetPathValue("agency_id", tt.agencyID)

			rec := httptest.NewRecorder()

			h.List(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}


func TestDeleteTourHandler(t *testing.T) {

	tests := []struct {
		name           string
		tourID         string
		mockUsecase    func(*mocks.MockDeleteTour)
		expectedStatus int
	}{
		{
			name:   "success",
			tourID: "550e8400-e29b-41d4-a716-446655440000",
			mockUsecase: func(m *mocks.MockDeleteTour) {
				m.ExecuteFunc = func(ctx context.Context, id uuid.UUID) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid uuid",
			tourID:         "bad-id",
			mockUsecase:    func(m *mocks.MockDeleteTour) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockDeleteTour{}
			tt.mockUsecase(mockUC)

			h := handler.NewTourHandler(nil, nil, nil, nil, mockUC)

			req := httptest.NewRequest(http.MethodDelete, "/tours/"+tt.tourID, nil)
			req.SetPathValue("tour_id", tt.tourID)

			rec := httptest.NewRecorder()

			h.Delete(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestUpdateTourHandler(t *testing.T) {

	tests := []struct {
		name           string
		tourID         string
		body           string
		mockUsecase    func(*mocks.MockUpdateTour)
		expectedStatus int
	}{
		{
			name:   "success",
			tourID: "550e8400-e29b-41d4-a716-446655440000",
			body: `{
				"agency_id":"550e8400-e29b-41d4-a716-446655440000",
				"name":"updated tour",
				"start_date":"2026-06-01T00:00:00Z",
				"end_date":"2026-06-10T00:00:00Z",
				"last_enrollment_date":"2026-05-25T00:00:00Z",
				"available_seat":20,
				"description":"this is an updated tour description",
				"price":2000,
				"discount":0
			}`,
			mockUsecase: func(m *mocks.MockUpdateTour) {
				m.ExecuteFunc = func(ctx context.Context, t *domain.Tour) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid uuid",
			tourID:         "bad-id",
			body:           `{}`,
			mockUsecase:    func(m *mocks.MockUpdateTour) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "invalid json",
			tourID:   "550e8400-e29b-41d4-a716-446655440000",
			body:     `{bad}`,
			mockUsecase: func(m *mocks.MockUpdateTour) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "usecase error",
			tourID: "550e8400-e29b-41d4-a716-446655440000",
			body: `{
				"agency_id":"550e8400-e29b-41d4-a716-446655440000",
				"name":"updated tour",
				"start_date":"2026-06-01T00:00:00Z",
				"end_date":"2026-06-10T00:00:00Z",
				"last_enrollment_date":"2026-05-25T00:00:00Z",
				"available_seat":20,
				"description":"this is an updated tour description",
				"price":2000,
				"discount":0
			}`,
			mockUsecase: func(m *mocks.MockUpdateTour) {
				m.ExecuteFunc = func(ctx context.Context, t *domain.Tour) error {
					return errors.New("update failed")
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockUpdateTour{}
			tt.mockUsecase(mockUC)

			h := handler.NewTourHandler(nil, nil, nil, mockUC, nil)

			req := httptest.NewRequest(http.MethodPut, "/tours/"+tt.tourID, bytes.NewBufferString(tt.body))
			req.SetPathValue("tour_id", tt.tourID)

			rec := httptest.NewRecorder()

			h.Update(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}