package bookingusecase_test

import (
	"context"
	"sort"

	"github.com/bishal05das/travelbuddy/internal/domain"
	mocks "github.com/bishal05das/travelbuddy/internal/mocks/repository"
	"github.com/google/uuid"
)

// inlineTx runs the function without a real transaction.
type inlineTx struct{}

func (inlineTx) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// fakeBookingRepo is an in-memory port.BookingRepository backed by the
// mock tour repository for tour details.
type fakeBookingRepo struct {
	tours    *mocks.MockTourRepository
	bookings map[uuid.UUID]*domain.BookingResponse
	created  []*domain.Booking
}

func newFakeBookingRepo(tours *mocks.MockTourRepository) *fakeBookingRepo {
	return &fakeBookingRepo{tours: tours, bookings: map[uuid.UUID]*domain.BookingResponse{}}
}

func (f *fakeBookingRepo) Create(ctx context.Context, b *domain.Booking) error {
	tour, err := f.tours.GetByID(ctx, b.TourID)
	if err != nil {
		return err
	}
	b.BookingID = uuid.New()
	f.created = append(f.created, b)
	f.bookings[b.BookingID] = &domain.BookingResponse{
		BookingID: b.BookingID, Status: b.Status, NumberOfPeople: b.NumberOfPeople,
		TotalPrice: b.TotalPrice, UserID: b.UserID, MemberID: b.MemberID,
		CustomerID: b.CustomerID, TourID: b.TourID, TourStartDate: tour.StartDate, AgencyID: tour.AgencyID,
	}
	return nil
}

func inScope(b *domain.BookingResponse, s domain.BookingScope) bool {
	if s.AgencyID != nil && b.AgencyID != *s.AgencyID {
		return false
	}
	if s.UserID != nil && (b.UserID == nil || *b.UserID != *s.UserID) {
		return false
	}
	return true
}

func (f *fakeBookingRepo) GetByID(_ context.Context, id uuid.UUID, s domain.BookingScope) (*domain.BookingResponse, error) {
	b, ok := f.bookings[id]
	if !ok || !inScope(b, s) {
		return nil, domain.ErrBookingNotFound
	}
	copied := *b
	return &copied, nil
}

func (f *fakeBookingRepo) List(_ context.Context, filter domain.BookingFilter) ([]*domain.BookingResponse, int, error) {
	var all []*domain.BookingResponse
	for _, b := range f.bookings {
		if inScope(b, filter.BookingScope) && (filter.Status == "" || b.Status == filter.Status) &&
			(filter.TourID == nil || b.TourID == *filter.TourID) {
			all = append(all, b)
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].BookingID.String() < all[j].BookingID.String() })
	start := (filter.Page - 1) * filter.Limit
	if start > len(all) {
		start = len(all)
	}
	end := min(start+filter.Limit, len(all))
	return all[start:end], len(all), nil
}

func (f *fakeBookingRepo) GetForUpdate(_ context.Context, id uuid.UUID, s domain.BookingScope) (*domain.LockedBooking, error) {
	b, ok := f.bookings[id]
	if !ok || !inScope(b, s) {
		return nil, domain.ErrBookingNotFound
	}
	return &domain.LockedBooking{
		BookingID: b.BookingID, TourID: b.TourID, AgencyID: b.AgencyID, UserID: b.UserID,
		NumberOfPeople: b.NumberOfPeople, Status: b.Status, TourStartDate: b.TourStartDate,
	}, nil
}

func (f *fakeBookingRepo) UpdateStatus(_ context.Context, id uuid.UUID, status string) error {
	b, ok := f.bookings[id]
	if !ok {
		return domain.ErrBookingNotFound
	}
	b.Status = status
	return nil
}

func (f *fakeBookingRepo) GetOrCreateCustomerByUser(context.Context, uuid.UUID) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (f *fakeBookingRepo) CreateCustomer(_ context.Context, c *domain.Customer) error {
	c.CustomerID = uuid.New()
	return nil
}

// fakePaymentRepo records payments and their status per booking.
type fakePaymentRepo struct {
	amounts []int
	status  map[uuid.UUID]string
}

func newFakePaymentRepo() *fakePaymentRepo {
	return &fakePaymentRepo{status: map[uuid.UUID]string{}}
}

func (f *fakePaymentRepo) Create(_ context.Context, p *domain.Payment) error {
	f.amounts = append(f.amounts, p.Amount)
	f.status[p.BookingID] = "pending"
	return nil
}

func (f *fakePaymentRepo) SetStatusForBooking(_ context.Context, bookingID uuid.UUID, from, to string) error {
	if f.status[bookingID] == from {
		f.status[bookingID] = to
	}
	return nil
}
