package tourusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/infrastructure/postgres/mocks"
	tourusecase "github.com/bishal05das/travelbuddy/internal/usecase/tour"
	"github.com/google/uuid"
)

func parseDate(t *testing.T, value string) time.Time {
	t.Helper()
	date, err := time.Parse("2006-01-02", value)
	if err != nil {
		t.Fatalf("failed to parse date %s: %v", value, err)
	}
	return date
}

func TestCreateTourUseCase(t *testing.T) {
	tests := map[string]struct {
		tour    *domain.Tour
		wantErr bool
	}{
		"successful create": {
			tour: &domain.Tour{
				AgencyID:           1,
				Name:               "Bandarban tour",
				StartDate:          parseDate(t, "2026-12-10"),
				EndDate:            parseDate(t, "2026-12-15"),
				Description:        "blah blah blah",
				LastEnrollmentDate: parseDate(t, "2026-12-05"),
				Price:              10000,
				Discount:           200,
			},
			wantErr: false,
		},
		"missing name": {
			tour: &domain.Tour{
				AgencyID:           1,
				Name:               "",
				StartDate:          parseDate(t, "2026-12-10"),
				EndDate:            parseDate(t, "2026-12-15"),
				Description:        "blah blah blah",
				LastEnrollmentDate: parseDate(t, "2026-12-05"),
				Price:              10000,
				Discount:           200,
			},
			wantErr: true,
		},
		"start>end": {
			tour: &domain.Tour{
				AgencyID:           1,
				Name:               "",
				StartDate:          parseDate(t, "2026-12-15"),
				EndDate:            parseDate(t, "2026-12-10"),
				Description:        "blah blah blah",
				LastEnrollmentDate: parseDate(t, "2026-12-05"),
				Price:              10000,
				Discount:           200,
			},
			wantErr: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			repo := mocks.NewMockTourRepository()
			uc := tourusecase.NewCreateTourUseCase(repo)
			err := uc.Execute(context.Background(),tt.tour)
			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.wantErr {
				if tt.tour.TourID == uuid.Nil {
					t.Error("expected tour ID to be set")
				}
			}
		})
	}
}
