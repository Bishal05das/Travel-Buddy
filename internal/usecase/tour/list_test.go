package tourusecase_test

import (
	"context"
	"testing"

	"github.com/bishal05das/travelbuddy/internal/domain"
	mocks "github.com/bishal05das/travelbuddy/internal/mocks/repository"
	tourusecase "github.com/bishal05das/travelbuddy/internal/usecase/tour"
	"github.com/google/uuid"
)

func TestListUseCase(t *testing.T) {
	agency1ID := uuid.New()
	agency2ID := uuid.New()
	agency3ID := uuid.New()
	tests := []struct {
		name          string
		seedTours     []*domain.Tour
		agencyID      uuid.UUID
		expectedCount int
	}{
		{
			name: "returns tours for given agency",
			seedTours: []*domain.Tour{
				{
					AgencyID:           agency1ID,
					Name:               "Bandarban tour",
					StartDate:          parseDate(t, "2026-12-10"),
					EndDate:            parseDate(t, "2026-12-15"),
					Description:        "blah blah blah",
					LastEnrollmentDate: parseDate(t, "2026-12-05"),
					Price:              10000,
					Discount:           200,
				},
				{
					AgencyID:           agency2ID,
					Name:               "sundarban tour",
					StartDate:          parseDate(t, "2026-12-11"),
					EndDate:            parseDate(t, "2026-12-15"),
					Description:        "blah blah blah",
					LastEnrollmentDate: parseDate(t, "2026-12-06"),
					Price:              10000,
					Discount:           200,
				},
				{
					AgencyID:           agency3ID,
					Name:               "sylhet tour",
					StartDate:          parseDate(t, "2026-12-10"),
					EndDate:            parseDate(t, "2026-12-15"),
					Description:        "blah blah blah",
					LastEnrollmentDate: parseDate(t, "2026-12-05"),
					Price:              10000,
					Discount:           200,
				},
			},
			agencyID:      agency3ID,
			expectedCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockTourRepository()

			for _, tour := range tt.seedTours {
				err := repo.CreateTour(context.Background(), tour)
				if err != nil {
					t.Fatalf("failed to seed tour: %v", err)
				}
			}

			usecase := tourusecase.NewListTourUseCase(repo)
			tours, err := usecase.Execute(context.Background(), tt.agencyID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tours) != tt.expectedCount {
				t.Fatalf("expected %d tours,got %d", tt.expectedCount, len(tours))
			}
		})
	}
}
