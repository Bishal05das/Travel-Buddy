package handler_test

import (
	"bytes"
	"context"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
	"github.com/bishal05das/travelbuddy/internal/domain"
	mocks "github.com/bishal05das/travelbuddy/internal/mocks/usecase"
	"github.com/google/uuid"
)

// pngBytes is a minimal payload that sniffs as image/png.
var pngBytes = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR")

// newAgencyForm builds a multipart request; an empty fileName omits the image part.
func newAgencyForm(t *testing.T, fields map[string]string, fileName string, content []byte) *http.Request {
	t.Helper()
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	for k, v := range fields {
		if err := mw.WriteField(k, v); err != nil {
			t.Fatal(err)
		}
	}
	if fileName != "" {
		fw, err := mw.CreateFormFile("image", fileName)
		if err != nil {
			t.Fatal(err)
		}
		fw.Write(content)
	}
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/agency", body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	return req
}

func TestCreateAgencyHandler(t *testing.T) {
	// The handler writes uploads relative to the working directory.
	t.Chdir(t.TempDir())

	validFields := map[string]string{
		"name":            "TravelPro",
		"address":         "Dhaka",
		"registration_id": "REG123",
	}

	tests := []struct {
		name           string
		request        func(t *testing.T) *http.Request
		mockUsecase    func(*mocks.MockCreateAgency)
		expectedStatus int
	}{
		{
			name: "success",
			request: func(t *testing.T) *http.Request {
				return newAgencyForm(t, validFields, "logo.png", pngBytes)
			},
			mockUsecase: func(m *mocks.MockCreateAgency) {
				m.ExecuteFunc = func(ctx context.Context, a *domain.Agency, imagePath string) error {
					if !strings.HasPrefix(imagePath, "images/agencies/") {
						return errors.New("unexpected image path " + imagePath)
					}
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "not a multipart form",
			request: func(t *testing.T) *http.Request {
				return httptest.NewRequest(http.MethodPost, "/agency", bytes.NewBufferString(`{bad-json}`))
			},
			mockUsecase:    func(m *mocks.MockCreateAgency) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing required fields",
			request: func(t *testing.T) *http.Request {
				return newAgencyForm(t, map[string]string{"name": "TravelPro"}, "logo.png", pngBytes)
			},
			mockUsecase:    func(m *mocks.MockCreateAgency) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "missing image",
			request: func(t *testing.T) *http.Request {
				return newAgencyForm(t, validFields, "", nil)
			},
			mockUsecase:    func(m *mocks.MockCreateAgency) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "unsupported image extension",
			request: func(t *testing.T) *http.Request {
				return newAgencyForm(t, validFields, "logo.gif", pngBytes)
			},
			mockUsecase:    func(m *mocks.MockCreateAgency) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "html disguised as png",
			request: func(t *testing.T) *http.Request {
				return newAgencyForm(t, validFields, "logo.png", []byte("<html><script>alert(1)</script></html>"))
			},
			mockUsecase:    func(m *mocks.MockCreateAgency) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "usecase error",
			request: func(t *testing.T) *http.Request {
				return newAgencyForm(t, validFields, "logo.png", pngBytes)
			},
			mockUsecase: func(m *mocks.MockCreateAgency) {
				m.ExecuteFunc = func(ctx context.Context, a *domain.Agency, imagePath string) error {
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

			h := handler.NewAgencyHandler(mockUC, nil, nil, nil)
			rec := httptest.NewRecorder()

			h.CreateAgency(rec, tt.request(t))

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected %d got %d: %s", tt.expectedStatus, rec.Code, rec.Body.String())
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

			h := handler.NewAgencyHandler(nil, mockUC, nil, nil)

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

			h := handler.NewAgencyHandler(nil, nil, mockUC, nil)

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
