package handler_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
	"github.com/bishal05das/travelbuddy/internal/domain"
	mocks "github.com/bishal05das/travelbuddy/internal/mocks/usecase"
	util "github.com/bishal05das/travelbuddy/utils"

	"github.com/google/uuid"
)

// helper to generate fake JWT token with payload
func generateTestJWT(payload util.Payload) string {
	payloadBytes, _ := json.Marshal(payload)
	encodedPayload := base64.URLEncoding.WithPadding(base64.NoPadding).
		EncodeToString(payloadBytes)
	return "header." + encodedPayload + ".signature"
}

func TestBookingHandler_CreateBooking(t *testing.T) {
	tourID := uuid.New()
	//userID := uuid.New()

	validReq := domain.BookingRequest{
		CustomerID:     uuid.New(),
		NumberOfPeople: 2,
		TotalPrice:     500,
		Status:         "pending",
		Method:         "cash",
		CustomerName:   "John Doe",
		CustomerEmail:  "john@example.com",
		CustomerPhone:  "+8801712345678",
	}

	validResponse := &domain.BookingResponse{
		BookingID:      uuid.New(),
		CustomerID:     validReq.CustomerID,
		TourID:         tourID,
		TourName:       "Test Tour",
		AgencyName:     "Test Agency",
		BookingDate:    time.Now(),
		NumberOfPeople: validReq.NumberOfPeople,
		TotalPrice:     float64(validReq.TotalPrice),
		Status:         validReq.Status,
		CreatedBy:      "user",
		CreatedAt:      time.Now(),
	}

	tests := []struct {
		name           string
		role           string
		tourID         string
		body           interface{}
		mockExecute    func(*mocks.MockCreateBookingUC)
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "invalid tour id",
			tourID:         "not-a-uuid",
			body:           validReq,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing authorization header",
			tourID:         tourID.String(),
			body:           validReq,
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid json body",
			tourID:         tourID.String(),
			role:           "user",
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "validation error",
			tourID:         tourID.String(),
			role:           "user",
			body:           domain.BookingRequest{CustomerName: ""},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "success booking user",
			tourID: tourID.String(),
			role:   "user",
			body:   validReq,
			mockExecute: func(m *mocks.MockCreateBookingUC) {
				m.ExecuteFunc = func(_ context.Context, _ *domain.BookingRequest, userID *uuid.UUID, memberID *uuid.UUID) (*domain.BookingResponse, error) {
					if userID == nil || memberID != nil {
						t.Fatal("wrong user/member assignment")
					}
					return validResponse, nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "success booking member",
			tourID: tourID.String(),
			role:   "member",
			body:   validReq,
			mockExecute: func(m *mocks.MockCreateBookingUC) {
				m.ExecuteFunc = func(_ context.Context, _ *domain.BookingRequest, userID *uuid.UUID, memberID *uuid.UUID) (*domain.BookingResponse, error) {
					if userID != nil || memberID == nil {
						t.Fatal("wrong user/member assignment")
					}
					return validResponse, nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid role",
			tourID:         tourID.String(),
			role:           "admin",
			body:           validReq,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:   "usecase returns error",
			tourID: tourID.String(),
			role:   "user",
			body:   validReq,
			mockExecute: func(m *mocks.MockCreateBookingUC) {
				m.ExecuteFunc = func(_ context.Context, _ *domain.BookingRequest, _ *uuid.UUID, _ *uuid.UUID) (*domain.BookingResponse, error) {
					return nil, errors.New("booking failed")
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockUC := &mocks.MockCreateBookingUC{}
			if tt.mockExecute != nil {
				tt.mockExecute(mockUC)
			}

			h := handler.NewBookingHandler(mockUC)

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			default:
				bodyBytes, _ = json.Marshal(v)
			}

			req := httptest.NewRequest(http.MethodPost, "/tours/"+tt.tourID+"/booking", bytes.NewBuffer(bodyBytes))

			// Set Authorization header if role is given
			if tt.role != "" {
				payload := util.Payload{
					UserID: uuid.New(),
					Role:   tt.role,
				}
				token := generateTestJWT(payload)
				req.Header.Set("Authorization", "Bearer "+token)
			} else if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			req.SetPathValue("tour_id", tt.tourID)

			rr := httptest.NewRecorder()

			h.CreateBooking(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d, body: %s", tt.expectedStatus, rr.Code, rr.Body.String())
			}
		})
	}
}
