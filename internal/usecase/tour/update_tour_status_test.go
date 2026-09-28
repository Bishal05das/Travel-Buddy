package tourusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	mocks "github.com/bishal05das/travelbuddy/internal/mocks/repository"
	bookingusecase "github.com/bishal05das/travelbuddy/internal/usecase/booking"
	tourusecase "github.com/bishal05das/travelbuddy/internal/usecase/tour"
	"github.com/google/uuid"
)

type cancelFixture struct {
	ctx      context.Context
	tours    *mocks.MockTourRepository
	bookings *mocks.MockBookingRepository
	payments *mocks.MockPaymentRepository
	tourID   uuid.UUID
	agency   uuid.UUID
	member   domain.Actor
	user     domain.Actor
	// pending (2 people), confirmed (3 people), already cancelled by the customer (1 person)
	pending, confirmed, selfCancelled uuid.UUID
}

func newCancelFixture(t *testing.T) *cancelFixture {
	t.Helper()
	f := &cancelFixture{ctx: context.Background(), agency: uuid.New()}
	f.member = domain.Actor{ID: uuid.New(), Role: domain.RoleMember, AgencyID: &f.agency}
	f.user = domain.Actor{ID: uuid.New(), Role: domain.RoleUser}

	f.tours = mocks.NewMockTourRepository()
	tour := &domain.Tour{
		AgencyID: f.agency, Name: "Sajek", TotalSeat: 10, Price: 1000,
		StartDate: time.Now().Add(72 * time.Hour), LastEnrollmentDate: time.Now().Add(24 * time.Hour),
	}
	if err := f.tours.CreateTour(f.ctx, tour); err != nil {
		t.Fatal(err)
	}
	f.tourID = tour.TourID
	f.bookings = mocks.NewMockBookingRepository(f.tours)
	f.payments = mocks.NewMockPaymentRepository()

	create := bookingusecase.NewCreateBookingUseCase(mocks.InlineTx{}, f.bookings, f.tours, f.payments)
	book := func(people int) uuid.UUID {
		resp, err := create.Execute(f.ctx, &domain.BookingCommand{TourID: f.tourID, UserID: &f.user.ID, NumberOfPeople: people, TotalPrice: 1000 * people})
		if err != nil {
			t.Fatal(err)
		}
		return resp.BookingID
	}
	f.pending, f.confirmed, f.selfCancelled = book(2), book(3), book(1)

	update := bookingusecase.NewUpdateBookingStatusUseCase(mocks.InlineTx{}, f.bookings, f.tours, f.payments)
	if _, err := update.Execute(f.ctx, f.member, f.agency, f.confirmed, domain.BookingConfirmed); err != nil {
		t.Fatal(err)
	}
	cancel := bookingusecase.NewCancelMyBookingUseCase(mocks.InlineTx{}, f.bookings, f.tours, f.payments)
	if _, err := cancel.Execute(f.ctx, f.user, f.selfCancelled); err != nil {
		t.Fatal(err)
	}
	f.expectSeats(t, 5) // 10 - 2 - 3
	return f
}

func (f *cancelFixture) expectSeats(t *testing.T, want int) {
	t.Helper()
	tour, _ := f.tours.GetByID(f.ctx, f.tourID)
	if tour.AvailableSeat != want {
		t.Fatalf("expected %d available seats, got %d", want, tour.AvailableSeat)
	}
}

func (f *cancelFixture) statusUC() interface {
	Execute(context.Context, domain.Actor, uuid.UUID, string) (*domain.TourStatusChange, error)
} {
	return tourusecase.NewUpdateTourStatusUseCase(mocks.InlineTx{}, f.tours, f.bookings, f.payments)
}

func TestCancellingTourCancelsItsBookings(t *testing.T) {
	f := newCancelFixture(t)

	change, err := f.statusUC().Execute(f.ctx, f.member, f.tourID, domain.TourCancelled)
	if err != nil {
		t.Fatal(err)
	}
	if change.CancelledBookings != 2 {
		t.Fatalf("expected 2 active bookings cancelled, got %d", change.CancelledBookings)
	}
	f.expectSeats(t, 10)

	// What the customer now sees for each booking.
	mine := bookingusecase.NewGetBookingUseCase(f.bookings)
	scope := domain.BookingScope{UserID: &f.user.ID}
	cases := []struct {
		id                   uuid.UUID
		reason, paymentAfter string
	}{
		{f.pending, domain.CancelledWithTour, "failed"},    // never verified
		{f.confirmed, domain.CancelledWithTour, "success"}, // verified: refunded offline
		{f.selfCancelled, domain.CancelledByCustomer, "failed"},
	}
	for _, c := range cases {
		b, err := mine.Execute(f.ctx, f.user, scope, c.id)
		if err != nil {
			t.Fatal(err)
		}
		if b.Status != domain.BookingCancelled || b.CancellationReason != c.reason || b.TourStatus != domain.TourCancelled {
			t.Fatalf("booking %s: status=%s reason=%s tour_status=%s", c.id, b.Status, b.CancellationReason, b.TourStatus)
		}
		if got := f.payments.Status[c.id]; got != c.paymentAfter {
			t.Fatalf("booking %s: expected payment %s, got %s", c.id, c.paymentAfter, got)
		}
	}

	// Nothing can be booked or reopened afterwards.
	create := bookingusecase.NewCreateBookingUseCase(mocks.InlineTx{}, f.bookings, f.tours, f.payments)
	if _, err := create.Execute(f.ctx, &domain.BookingCommand{TourID: f.tourID, UserID: &f.user.ID, NumberOfPeople: 1, TotalPrice: 1000}); err == nil {
		t.Fatal("expected booking a cancelled tour to fail")
	}
	if _, err := f.statusUC().Execute(f.ctx, f.member, f.tourID, domain.TourOpen); !errors.Is(err, domain.ErrTourCancelled) {
		t.Fatalf("expected reopening to fail with ErrTourCancelled, got %v", err)
	}
	if _, err := f.statusUC().Execute(f.ctx, f.member, f.tourID, domain.TourCancelled); !errors.Is(err, domain.ErrTourCancelled) {
		t.Fatalf("expected a second cancellation to fail, got %v", err)
	}
	f.expectSeats(t, 10)
}

func TestClosingTourKeepsBookings(t *testing.T) {
	f := newCancelFixture(t)

	change, err := f.statusUC().Execute(f.ctx, f.member, f.tourID, domain.TourClosed)
	if err != nil {
		t.Fatal(err)
	}
	if change.CancelledBookings != 0 {
		t.Fatalf("closing must not cancel bookings, cancelled %d", change.CancelledBookings)
	}
	b, _ := f.bookings.GetByID(f.ctx, f.confirmed, domain.BookingScope{})
	if b.Status != domain.BookingConfirmed {
		t.Fatalf("expected booking to stay confirmed, got %s", b.Status)
	}
	f.expectSeats(t, 5)

	if _, err := f.statusUC().Execute(f.ctx, f.member, f.tourID, domain.TourOpen); err != nil {
		t.Fatalf("closed tours can be reopened: %v", err)
	}
}

func TestOtherAgencyCannotCancelTour(t *testing.T) {
	f := newCancelFixture(t)
	other := uuid.New()
	outsider := domain.Actor{ID: uuid.New(), Role: domain.RoleMember, AgencyID: &other}

	if _, err := f.statusUC().Execute(f.ctx, outsider, f.tourID, domain.TourCancelled); err == nil {
		t.Fatal("expected another agency to be refused")
	}
	b, _ := f.bookings.GetByID(f.ctx, f.pending, domain.BookingScope{})
	if b.Status != domain.BookingPending {
		t.Fatalf("bookings must be untouched, got %s", b.Status)
	}
	f.expectSeats(t, 5)
}

func TestInvalidTourStatus(t *testing.T) {
	f := newCancelFixture(t)
	if _, err := f.statusUC().Execute(f.ctx, f.member, f.tourID, "archived"); err == nil {
		t.Fatal("expected an unknown status to be rejected")
	}
}
