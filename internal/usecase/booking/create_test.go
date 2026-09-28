package bookingusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	mocks "github.com/bishal05das/travelbuddy/internal/mocks/repository"
	bookingusecase "github.com/bishal05das/travelbuddy/internal/usecase/booking"
	"github.com/google/uuid"
)

// inlineTx runs the function without a real transaction.
type inlineTx struct{}

func (inlineTx) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type fakeBookingRepo struct{ created []*domain.Booking }

func (f *fakeBookingRepo) Create(_ context.Context, b *domain.Booking) error {
	b.BookingID = uuid.New()
	f.created = append(f.created, b)
	return nil
}
func (f *fakeBookingRepo) GetByID(_ context.Context, id uuid.UUID) (*domain.BookingResponse, error) {
	return &domain.BookingResponse{BookingID: id}, nil
}
func (f *fakeBookingRepo) Update(context.Context, *domain.Booking) error { return nil }
func (f *fakeBookingRepo) Cancel(context.Context, uuid.UUID) error       { return nil }
func (f *fakeBookingRepo) GetOrCreateCustomerByUser(context.Context, uuid.UUID) (uuid.UUID, error) {
	return uuid.New(), nil
}
func (f *fakeBookingRepo) CreateCustomer(_ context.Context, c *domain.Customer) error {
	c.CustomerID = uuid.New()
	return nil
}

type fakePaymentRepo struct{ amounts []int }

func (f *fakePaymentRepo) Create(_ context.Context, p *domain.Payment) error {
	f.amounts = append(f.amounts, p.Amount)
	return nil
}

func TestCreateBooking(t *testing.T) {
	agencyID := uuid.New()
	otherAgency := uuid.New()
	userID := uuid.New()
	memberID := uuid.New()
	guest := &domain.GuestCustomer{Name: "Guest", Email: "g@test.com", Phone: "+8801700000000"}

	// Price 1000 with 10% discount -> 900 per person.
	newTour := func() *domain.Tour {
		return &domain.Tour{
			TourID:             uuid.New(),
			AgencyID:           agencyID,
			Name:               "Sajek",
			Status:             "open",
			AvailableSeat:      5,
			Price:              1000,
			Discount:           10,
			LastEnrollmentDate: time.Now().Add(24 * time.Hour),
		}
	}

	tests := []struct {
		name      string
		mutate    func(*domain.Tour)
		cmd       func(tourID uuid.UUID) *domain.BookingCommand
		wantErr   bool
		wantSeats int
	}{
		{
			name: "user books 3 people at the group price",
			cmd: func(id uuid.UUID) *domain.BookingCommand {
				return &domain.BookingCommand{TourID: id, UserID: &userID, NumberOfPeople: 3, TotalPrice: 2700}
			},
			wantSeats: 2,
		},
		{
			name: "booking exactly the remaining seats is allowed",
			cmd: func(id uuid.UUID) *domain.BookingCommand {
				return &domain.BookingCommand{TourID: id, UserID: &userID, NumberOfPeople: 5, TotalPrice: 4500}
			},
			wantSeats: 0,
		},
		{
			name: "single-person price for a group is rejected",
			cmd: func(id uuid.UUID) *domain.BookingCommand {
				return &domain.BookingCommand{TourID: id, UserID: &userID, NumberOfPeople: 3, TotalPrice: 900}
			},
			wantErr: true, wantSeats: 5,
		},
		{
			name: "more people than seats",
			cmd: func(id uuid.UUID) *domain.BookingCommand {
				return &domain.BookingCommand{TourID: id, UserID: &userID, NumberOfPeople: 6, TotalPrice: 5400}
			},
			wantErr: true, wantSeats: 5,
		},
		{
			name:   "closed tour",
			mutate: func(t *domain.Tour) { t.Status = "closed" },
			cmd: func(id uuid.UUID) *domain.BookingCommand {
				return &domain.BookingCommand{TourID: id, UserID: &userID, NumberOfPeople: 1, TotalPrice: 900}
			},
			wantErr: true, wantSeats: 5,
		},
		{
			name:   "enrollment deadline passed",
			mutate: func(t *domain.Tour) { t.LastEnrollmentDate = time.Now().Add(-time.Hour) },
			cmd: func(id uuid.UUID) *domain.BookingCommand {
				return &domain.BookingCommand{TourID: id, UserID: &userID, NumberOfPeople: 1, TotalPrice: 900}
			},
			wantErr: true, wantSeats: 5,
		},
		{
			name: "member books a guest on their own agency's tour",
			cmd: func(id uuid.UUID) *domain.BookingCommand {
				return &domain.BookingCommand{TourID: id, MemberID: &memberID, AgencyID: &agencyID, GuestInfo: guest, NumberOfPeople: 2, TotalPrice: 1800}
			},
			wantSeats: 3,
		},
		{
			name: "member cannot book another agency's tour",
			cmd: func(id uuid.UUID) *domain.BookingCommand {
				return &domain.BookingCommand{TourID: id, MemberID: &memberID, AgencyID: &otherAgency, GuestInfo: guest, NumberOfPeople: 1, TotalPrice: 900}
			},
			wantErr: true, wantSeats: 5,
		},
		{
			name: "member booking without guest details",
			cmd: func(id uuid.UUID) *domain.BookingCommand {
				return &domain.BookingCommand{TourID: id, MemberID: &memberID, AgencyID: &agencyID, NumberOfPeople: 1, TotalPrice: 900}
			},
			wantErr: true, wantSeats: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tourRepo := mocks.NewMockTourRepository()
			tour := newTour()
			if tt.mutate != nil {
				tt.mutate(tour)
			}
			if err := tourRepo.CreateTour(context.Background(), tour); err != nil {
				t.Fatal(err)
			}
			bookings := &fakeBookingRepo{}
			payments := &fakePaymentRepo{}
			uc := bookingusecase.NewCreateBookingUseCase(inlineTx{}, bookings, tourRepo, payments)

			cmd := tt.cmd(tour.TourID)
			_, err := uc.Execute(context.Background(), cmd)
			if tt.wantErr != (err != nil) {
				t.Fatalf("wantErr=%v, got %v", tt.wantErr, err)
			}

			got, _ := tourRepo.GetByID(context.Background(), tour.TourID)
			if got.AvailableSeat != tt.wantSeats {
				t.Fatalf("expected %d seats left, got %d", tt.wantSeats, got.AvailableSeat)
			}
			if !tt.wantErr && (len(payments.amounts) != 1 || payments.amounts[0] != cmd.TotalPrice) {
				t.Fatalf("expected one payment of %d, got %v", cmd.TotalPrice, payments.amounts)
			}
		})
	}
}
