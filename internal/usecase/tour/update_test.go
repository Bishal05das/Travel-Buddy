package tourusecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bishal05das/travelbuddy/internal/domain"
	mocks "github.com/bishal05das/travelbuddy/internal/mocks/repository"
	tourusecase "github.com/bishal05das/travelbuddy/internal/usecase/tour"
	"github.com/google/uuid"
)

func TestUpdateTourUseCase(t *testing.T) {

	tour1ID := uuid.New()
	tour2ID := uuid.New()
	tour3ID := uuid.New()
	agency1ID := uuid.New()
	agency2ID := uuid.New()

	tests := []struct {
		name        string
		seedTours   []*domain.Tour
		updateTour  *domain.Tour
		expectedErr bool
	}{
		{
			name: "successfully updates existing tour",
			seedTours: []*domain.Tour{
				{
					TourID:             tour1ID,
					AgencyID:           agency1ID,
					Name:               "Bandarban tour",
					StartDate:          parseDate(t, "2026-12-10"),
					EndDate:            parseDate(t, "2026-12-15"),
					Description:        "blah blah blah",
					LastEnrollmentDate: parseDate(t, "2026-12-05"),
					Price:              10000,
					Discount:           20,
					TotalSeat:          30,
				},
				{
					TourID:             tour2ID,
					AgencyID:           agency1ID,
					Name:               "sundarban tour",
					StartDate:          parseDate(t, "2026-12-11"),
					EndDate:            parseDate(t, "2026-12-15"),
					Description:        "blah blah blah",
					LastEnrollmentDate: parseDate(t, "2026-12-06"),
					Price:              10000,
					Discount:           20,
					TotalSeat:          30,
				},
				{
					TourID:             tour3ID,
					AgencyID:           agency2ID,
					Name:               "sylhet tour",
					StartDate:          parseDate(t, "2026-12-10"),
					EndDate:            parseDate(t, "2026-12-15"),
					Description:        "blah blah blah",
					LastEnrollmentDate: parseDate(t, "2026-12-05"),
					Price:              10000,
					Discount:           20,
					TotalSeat:          30,
				},
			},
			updateTour: &domain.Tour{
				TourID:             tour3ID,
				AgencyID:           agency2ID,
				Name:               "sylhet Tour",
				StartDate:          parseDate(t, "2026-12-10"),
				EndDate:            parseDate(t, "2026-12-15"),
				Description:        "blah blah blah",
				LastEnrollmentDate: parseDate(t, "2026-12-05"),
				Price:              10000,
				Discount:           20,
				TotalSeat:          30,
			},
			expectedErr: false,
		},
		{
			name:        "fails when tour does not exist",
			seedTours:   []*domain.Tour{},
			updateTour:  &domain.Tour{TourID: uuid.New()},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockTourRepository()

			for _, tour := range tt.seedTours {
				err := repo.CreateTour(context.Background(), tour)
				if err != nil {
					t.Fatalf("failed to seed tour: %v", tour)
				}
			}
			usecase := tourusecase.NewUpdateTourUseCase(repo)
			err := usecase.Execute(context.Background(), tt.updateTour)
			if tt.expectedErr && err == nil {
				t.Fatalf("expected error,got nil")
			}
			if !tt.expectedErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestUpdateTourKeepsBookedSeats(t *testing.T) {
	agencyID := uuid.New()
	newTour := func() *domain.Tour {
		// Capacity 10 with 4 seats already booked.
		return &domain.Tour{
			TourID:             uuid.New(),
			AgencyID:           agencyID,
			Name:               "Sajek tour",
			StartDate:          parseDate(t, "2026-12-10"),
			EndDate:            parseDate(t, "2026-12-15"),
			Description:        "three days in the hills",
			LastEnrollmentDate: parseDate(t, "2026-12-05"),
			Price:              10000,
			TotalSeat:          10,
			AvailableSeat:      6,
		}
	}

	tests := []struct {
		name          string
		newTotal      int
		wantErr       error
		wantAvailable int
	}{
		{name: "raising capacity adds free seats", newTotal: 15, wantAvailable: 11},
		{name: "lowering capacity keeps bookings", newTotal: 5, wantAvailable: 1},
		{name: "capacity equal to booked sells out", newTotal: 4, wantAvailable: 0},
		{name: "capacity below booked is rejected", newTotal: 3, wantErr: domain.ErrCapacityBelowBooked, wantAvailable: 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockTourRepository()
			seed := newTour()
			if err := repo.CreateTour(context.Background(), seed); err != nil {
				t.Fatal(err)
			}

			update := newTour()
			update.TourID = seed.TourID
			update.TotalSeat = tt.newTotal
			update.AvailableSeat = 0 // ignored: remaining seats are derived

			err := tourusecase.NewUpdateTourUseCase(repo).Execute(context.Background(), update)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}

			got, _ := repo.GetByID(context.Background(), seed.TourID)
			if got.AvailableSeat != tt.wantAvailable {
				t.Fatalf("expected %d available seats, got %d", tt.wantAvailable, got.AvailableSeat)
			}
			if got.BookedSeats() != 4 {
				t.Fatalf("booked seats changed: %d", got.BookedSeats())
			}
		})
	}
}
