package tourusecase_test

import (
	"testing"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/infrastructure/postgres/mocks"
	tourusecase "github.com/bishal05das/travelbuddy/internal/usecase/tour"
)

func TestUpdateTourUseCase(t *testing.T) {
	tests := []struct {
		name        string
		seedTours   []*domain.Tour
		updateTour  *domain.Tour
		expectedErr bool
	}{
		{
			name: "successfully upodates existing tour",
			seedTours: []*domain.Tour{
				{
					AgencyID:           1,
					Name:               "Bandarban tour",
					StartDate:          parseDate(t, "2026-12-10"),
					EndDate:            parseDate(t, "2026-12-15"),
					Description:        "blah blah blah",
					LastEnrollmentDate: parseDate(t, "2026-12-05"),
					Price:              10000,
					Discount:           200,
				},
				{
					AgencyID:           1,
					Name:               "sundarban tour",
					StartDate:          parseDate(t, "2026-12-11"),
					EndDate:            parseDate(t, "2026-12-15"),
					Description:        "blah blah blah",
					LastEnrollmentDate: parseDate(t, "2026-12-06"),
					Price:              10000,
					Discount:           200,
				},
				{
					AgencyID:           2,
					Name:               "sylhet tour",
					StartDate:          parseDate(t, "2026-12-10"),
					EndDate:            parseDate(t, "2026-12-15"),
					Description:        "blah blah blah",
					LastEnrollmentDate: parseDate(t, "2026-12-05"),
					Price:              10000,
					Discount:           200,
				},
			},
			updateTour: &domain.Tour{
				ID:                 3,
				AgencyID:           2,
				Name:               "sylhet Tour",
				StartDate:          parseDate(t, "2026-12-10"),
				EndDate:            parseDate(t, "2026-12-15"),
				Description:        "blah blah blah",
				LastEnrollmentDate: parseDate(t, "2026-12-05"),
				Price:              10000,
				Discount:           200,
			},
			expectedErr: false,
		},
		{
			name:       "fails when tour does not exist",
			seedTours:  []*domain.Tour{},
			updateTour: &domain.Tour{ID: 99},
			expectedErr: true,
		},
	}

	for _,tt := range tests{
		t.Run(tt.name,func(t *testing.T){
			repo := mocks.NewMockTourRepository()

			for _,tour := range tt.seedTours {
				err := repo.CreateTour(tour)
				if err != nil {
					t.Fatalf("failed to seed tour: %v",tour)
				}
			}
			usecase := tourusecase.NewUpdateTourUseCase(repo)
			err := usecase.Execute(tt.updateTour)
			if tt.expectedErr && err == nil {
				t.Fatalf("expected error,got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Fatalf("unexpected error: %v",err)
			}
		})
	}
}
