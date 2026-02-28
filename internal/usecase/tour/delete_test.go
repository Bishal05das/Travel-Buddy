package tourusecase_test

import (
	"context"
	"testing"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/infrastructure/postgres/mocks"
	tourusecase "github.com/bishal05das/travelbuddy/internal/usecase/tour"
	"github.com/google/uuid"
)

func TestDeleteTourUseCase(t *testing.T) {
	tour1ID := uuid.New()
	tour2ID := uuid.New()
	tour3ID := uuid.New()
	agency1ID := uuid.New()
	agency2ID := uuid.New()
	tests := []struct {
		name      string
		seedTours []*domain.Tour
		tourID    uuid.UUID
		expectErr bool
	}{
		{
			name: "successful delete",
			seedTours: []*domain.Tour{
				{
					TourID: tour1ID,
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
					TourID: tour2ID,
					AgencyID:           agency1ID,
					Name:               "sundarban tour",
					StartDate:          parseDate(t, "2026-12-11"),
					EndDate:            parseDate(t, "2026-12-15"),
					Description:        "blah blah blah",
					LastEnrollmentDate: parseDate(t, "2026-12-06"),
					Price:              10000,
					Discount:           200,
				},
				{
					TourID: tour3ID,
					AgencyID:           agency2ID,
					Name:               "sylhet tour",
					StartDate:          parseDate(t, "2026-12-10"),
					EndDate:            parseDate(t, "2026-12-15"),
					Description:        "blah blah blah",
					LastEnrollmentDate: parseDate(t, "2026-12-05"),
					Price:              10000,
					Discount:           200,
				},
			},
			tourID:    tour1ID,
			expectErr: false,
		},
		{
			name:      "fails when tour does not exist",
			seedTours: []*domain.Tour{},
			tourID:    tour1ID,
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockTourRepository()

			for _, tour := range tt.seedTours {
				err := repo.CreateTour(context.Background(),tour)
				if err != nil {
					t.Fatalf("failed to seed tour: %v", err)
				}
			}
			usecase := tourusecase.NewDeleteTourUseCase(repo)
			err := usecase.Execute(context.Background(),tt.tourID)
			if tt.expectErr && err == nil {
				t.Fatalf("expected error,got nil")
			}
			if !tt.expectErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
