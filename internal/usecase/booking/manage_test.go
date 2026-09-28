package bookingusecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	mocks "github.com/bishal05das/travelbuddy/internal/mocks/repository"
	bookingusecase "github.com/bishal05das/travelbuddy/internal/usecase/booking"
	"github.com/google/uuid"
)

// fixture: agency A's tour (capacity 10, 900/person) with a 3-person user
// booking and a 2-person guest booking by an agency A member.
type fixture struct {
	ctx                  context.Context
	tours                *mocks.MockTourRepository
	bookings             *mocks.MockBookingRepository
	payments             *mocks.MockPaymentRepository
	tourID               uuid.UUID
	agencyA, agencyB     uuid.UUID
	userID, otherUserID  uuid.UUID
	memberA, memberB     domain.Actor
	user, otherUser      domain.Actor
	userBooking, guestBk uuid.UUID
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	f := &fixture{ctx: context.Background(), agencyA: uuid.New(), agencyB: uuid.New(), userID: uuid.New(), otherUserID: uuid.New()}
	f.memberA = domain.Actor{ID: uuid.New(), Role: domain.RoleMember, AgencyID: &f.agencyA}
	f.memberB = domain.Actor{ID: uuid.New(), Role: domain.RoleMember, AgencyID: &f.agencyB}
	f.user = domain.Actor{ID: f.userID, Role: domain.RoleUser}
	f.otherUser = domain.Actor{ID: f.otherUserID, Role: domain.RoleUser}

	f.tours = mocks.NewMockTourRepository()
	tour := &domain.Tour{
		AgencyID: f.agencyA, Name: "Sajek", Status: "open", TotalSeat: 10, Price: 1000, Discount: 10,
		StartDate: time.Now().Add(72 * time.Hour), LastEnrollmentDate: time.Now().Add(24 * time.Hour),
	}
	if err := f.tours.CreateTour(f.ctx, tour); err != nil {
		t.Fatal(err)
	}
	f.tourID = tour.TourID
	f.bookings = mocks.NewMockBookingRepository(f.tours)
	f.payments = mocks.NewMockPaymentRepository()

	create := bookingusecase.NewCreateBookingUseCase(mocks.InlineTx{}, f.bookings, f.tours, f.payments)
	userResp, err := create.Execute(f.ctx, &domain.BookingCommand{TourID: f.tourID, UserID: &f.userID, NumberOfPeople: 3, TotalPrice: 2700})
	if err != nil {
		t.Fatal(err)
	}
	guestResp, err := create.Execute(f.ctx, &domain.BookingCommand{
		TourID: f.tourID, MemberID: &f.memberA.ID, AgencyID: &f.agencyA, NumberOfPeople: 2, TotalPrice: 1800,
		GuestInfo: &domain.GuestCustomer{Name: "Guest", Email: "g@t.com", Phone: "+8801700000000"},
	})
	if err != nil {
		t.Fatal(err)
	}
	f.userBooking, f.guestBk = userResp.BookingID, guestResp.BookingID
	f.expectSeats(t, 5)
	return f
}

func (f *fixture) expectSeats(t *testing.T, want int) {
	t.Helper()
	tour, _ := f.tours.GetByID(f.ctx, f.tourID)
	if tour.AvailableSeat != want {
		t.Fatalf("expected %d available seats, got %d", want, tour.AvailableSeat)
	}
}

func (f *fixture) updater() interface {
	Execute(context.Context, domain.Actor, uuid.UUID, uuid.UUID, string) (*domain.BookingResponse, error)
} {
	return bookingusecase.NewUpdateBookingStatusUseCase(mocks.InlineTx{}, f.bookings, f.tours, f.payments)
}

func TestAgencyConfirmThenCancel(t *testing.T) {
	f := newFixture(t)

	resp, err := f.updater().Execute(f.ctx, f.memberA, f.agencyA, f.userBooking, domain.BookingConfirmed)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != domain.BookingConfirmed || f.payments.Status[f.userBooking] != "success" {
		t.Fatalf("confirm: status %s, payment %s", resp.Status, f.payments.Status[f.userBooking])
	}
	f.expectSeats(t, 5)

	if _, err := f.updater().Execute(f.ctx, f.memberA, f.agencyA, f.userBooking, domain.BookingCancelled); err != nil {
		t.Fatal(err)
	}
	f.expectSeats(t, 8) // 3 seats returned
	if f.payments.Status[f.userBooking] != "success" {
		t.Fatalf("a verified payment must stay success for an offline refund, got %s", f.payments.Status[f.userBooking])
	}

	// Final state: no further changes, and seats are not returned twice.
	_, err = f.updater().Execute(f.ctx, f.memberA, f.agencyA, f.userBooking, domain.BookingCancelled)
	if !errors.Is(err, domain.ErrInvalidBookingTransition) {
		t.Fatalf("expected invalid transition, got %v", err)
	}
	f.expectSeats(t, 8)
}

func TestAgencyStatusRules(t *testing.T) {
	f := newFixture(t)

	if _, err := f.updater().Execute(f.ctx, f.memberA, f.agencyA, f.guestBk, domain.BookingCompleted); !errors.Is(err, domain.ErrInvalidBookingTransition) {
		t.Fatalf("pending -> completed: expected invalid transition, got %v", err)
	}
	if _, err := f.updater().Execute(f.ctx, f.memberA, f.agencyA, f.guestBk, domain.BookingCancelled); err != nil {
		t.Fatal(err)
	}
	if f.payments.Status[f.guestBk] != "failed" {
		t.Fatalf("cancelling an unverified booking should fail its payment, got %s", f.payments.Status[f.guestBk])
	}
	f.expectSeats(t, 7)
}

func TestAgencyCannotTouchOtherAgencyBookings(t *testing.T) {
	f := newFixture(t)

	// Member B acting through their own agency cannot find agency A's booking.
	if _, err := f.updater().Execute(f.ctx, f.memberB, f.agencyB, f.userBooking, domain.BookingCancelled); err == nil {
		t.Fatal("expected agency B to be unable to cancel agency A's booking")
	}
	// Member B naming agency A is refused before any lookup.
	if _, err := f.updater().Execute(f.ctx, f.memberB, f.agencyA, f.userBooking, domain.BookingCancelled); err == nil {
		t.Fatal("expected agency B to be refused agency A's scope")
	}
	f.expectSeats(t, 5)
}

func TestUserCancelsOwnBooking(t *testing.T) {
	f := newFixture(t)
	cancel := bookingusecase.NewCancelMyBookingUseCase(mocks.InlineTx{}, f.bookings, f.tours, f.payments)

	if _, err := cancel.Execute(f.ctx, f.otherUser, f.userBooking); err == nil {
		t.Fatal("expected another user to be unable to cancel this booking")
	}
	f.expectSeats(t, 5)

	resp, err := cancel.Execute(f.ctx, f.user, f.userBooking)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != domain.BookingCancelled || f.payments.Status[f.userBooking] != "failed" {
		t.Fatalf("cancel: status %s, payment %s", resp.Status, f.payments.Status[f.userBooking])
	}
	f.expectSeats(t, 8)
}

func TestUserCannotCancelAfterTourStarts(t *testing.T) {
	f := newFixture(t)
	f.bookings.Bookings[f.userBooking].TourStartDate = time.Now().Add(-time.Hour)
	cancel := bookingusecase.NewCancelMyBookingUseCase(mocks.InlineTx{}, f.bookings, f.tours, f.payments)

	if _, err := cancel.Execute(f.ctx, f.user, f.userBooking); !errors.Is(err, bookingusecase.ErrTourAlreadyStarted) {
		t.Fatalf("expected ErrTourAlreadyStarted, got %v", err)
	}
	f.expectSeats(t, 5)
}

func TestListBookingsScopes(t *testing.T) {
	f := newFixture(t)
	list := bookingusecase.NewListBookingsUseCase(f.bookings)

	count := func(actor domain.Actor, filter domain.BookingFilter) (int, error) {
		page, err := list.Execute(f.ctx, actor, filter)
		if err != nil {
			return 0, err
		}
		return page.Meta.TotalCount, nil
	}

	if n, err := count(f.memberA, domain.BookingFilter{BookingScope: domain.BookingScope{AgencyID: &f.agencyA}}); err != nil || n != 2 {
		t.Fatalf("agency A: expected 2 bookings, got %d (%v)", n, err)
	}
	if n, err := count(f.memberA, domain.BookingFilter{BookingScope: domain.BookingScope{AgencyID: &f.agencyA}, Status: domain.BookingConfirmed}); err != nil || n != 0 {
		t.Fatalf("status filter: expected 0, got %d (%v)", n, err)
	}
	if n, err := count(f.user, domain.BookingFilter{BookingScope: domain.BookingScope{UserID: &f.userID}}); err != nil || n != 1 {
		t.Fatalf("user: expected 1 booking, got %d (%v)", n, err)
	}
	if n, err := count(f.otherUser, domain.BookingFilter{BookingScope: domain.BookingScope{UserID: &f.otherUserID}}); err != nil || n != 0 {
		t.Fatalf("other user: expected 0 bookings, got %d (%v)", n, err)
	}

	forbidden := []struct {
		name   string
		actor  domain.Actor
		filter domain.BookingFilter
	}{
		{"member B listing agency A", f.memberB, domain.BookingFilter{BookingScope: domain.BookingScope{AgencyID: &f.agencyA}}},
		{"user listing another user", f.otherUser, domain.BookingFilter{BookingScope: domain.BookingScope{UserID: &f.userID}}},
		{"user listing an agency", f.user, domain.BookingFilter{BookingScope: domain.BookingScope{AgencyID: &f.agencyA}}},
		{"member without scope", f.memberA, domain.BookingFilter{}},
	}
	for _, tt := range forbidden {
		if _, err := count(tt.actor, tt.filter); err == nil {
			t.Fatalf("%s: expected an error", tt.name)
		}
	}
	if _, err := count(f.memberA, domain.BookingFilter{BookingScope: domain.BookingScope{AgencyID: &f.agencyA}, Status: "lost"}); err == nil {
		t.Fatal("expected an invalid status filter to be rejected")
	}
}

func TestCheckBookingTransition(t *testing.T) {
	allowed := [][2]string{{"pending", "confirmed"}, {"pending", "cancelled"}, {"confirmed", "cancelled"}, {"confirmed", "completed"}}
	for _, tr := range allowed {
		if err := domain.CheckBookingTransition(tr[0], tr[1]); err != nil {
			t.Errorf("%s -> %s should be allowed: %v", tr[0], tr[1], err)
		}
	}
	denied := [][2]string{{"pending", "completed"}, {"cancelled", "confirmed"}, {"completed", "cancelled"}, {"confirmed", "pending"}}
	for _, tr := range denied {
		if err := domain.CheckBookingTransition(tr[0], tr[1]); err == nil {
			t.Errorf("%s -> %s should be denied", tr[0], tr[1])
		}
	}
}
