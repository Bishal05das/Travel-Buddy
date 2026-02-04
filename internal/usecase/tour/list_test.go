package tourusecase_test

import (
	"testing"
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/repository/mocks"
	tourusecase "github.com/bishal05das/travelbuddy/internal/usecase/tour"
)


func TestListUseCase(t *testing.T) {
	tests := []struct {
		name          string
		seedTours     []*domain.Tour
		agencyID      int
		expectedCount int
	}{
		{
			name: "returns tours for given agency",
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
			agencyID:      2,
			expectedCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockTourRepository()

			for _, tour := range tt.seedTours {
				err := repo.CreateTour(tour)
				if err != nil {
					t.Fatalf("failed to sed tour: %v", err)
				}
			}

			usecase := tourusecase.NewListTourUseCase(repo)
			tours, err := usecase.Execute(tt.agencyID)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tours) != tt.expectedCount {
				t.Fatalf("expected %d tours,got %d", tt.expectedCount, len(tours))
			}
		})
	}
}
