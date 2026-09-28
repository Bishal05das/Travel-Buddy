package mocks

import (
	"context"
	"sort"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

var (
	_ port.BookingRepository = (*MockBookingRepository)(nil)
	_ port.PaymentRepository = (*MockPaymentRepository)(nil)
	_ port.TxManager         = InlineTx{}
)

// InlineTx runs the function without a real transaction.
type InlineTx struct{}

func (InlineTx) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// MockBookingRepository is an in-memory port.BookingRepository that reads
// tour details from a MockTourRepository.
type MockBookingRepository struct {
	tours    *MockTourRepository
	Bookings map[uuid.UUID]*domain.BookingResponse
	Created  []*domain.Booking
}

func NewMockBookingRepository(tours *MockTourRepository) *MockBookingRepository {
	return &MockBookingRepository{tours: tours, Bookings: map[uuid.UUID]*domain.BookingResponse{}}
}

func (m *MockBookingRepository) Create(ctx context.Context, b *domain.Booking) error {
	tour, err := m.tours.GetByID(ctx, b.TourID)
	if err != nil {
		return err
	}
	b.BookingID = uuid.New()
	m.Created = append(m.Created, b)
	m.Bookings[b.BookingID] = &domain.BookingResponse{
		BookingID: b.BookingID, Status: b.Status, NumberOfPeople: b.NumberOfPeople,
		TotalPrice: b.TotalPrice, UserID: b.UserID, MemberID: b.MemberID,
		CustomerID: b.CustomerID, TourID: b.TourID, TourStartDate: tour.StartDate, AgencyID: tour.AgencyID,
	}
	return nil
}

func bookingInScope(b *domain.BookingResponse, s domain.BookingScope) bool {
	if s.AgencyID != nil && b.AgencyID != *s.AgencyID {
		return false
	}
	if s.UserID != nil && (b.UserID == nil || *b.UserID != *s.UserID) {
		return false
	}
	return true
}

// view returns a copy with the tour's current status, like the SQL join.
func (m *MockBookingRepository) view(b *domain.BookingResponse) *domain.BookingResponse {
	copied := *b
	if tour, err := m.tours.GetByID(context.Background(), b.TourID); err == nil {
		copied.TourStatus = tour.Status
	}
	return &copied
}

func (m *MockBookingRepository) GetByID(_ context.Context, id uuid.UUID, s domain.BookingScope) (*domain.BookingResponse, error) {
	b, ok := m.Bookings[id]
	if !ok || !bookingInScope(b, s) {
		return nil, domain.ErrBookingNotFound
	}
	return m.view(b), nil
}

func (m *MockBookingRepository) List(_ context.Context, f domain.BookingFilter) ([]*domain.BookingResponse, int, error) {
	var all []*domain.BookingResponse
	for _, b := range m.Bookings {
		if bookingInScope(b, f.BookingScope) && (f.Status == "" || b.Status == f.Status) &&
			(f.TourID == nil || b.TourID == *f.TourID) {
			all = append(all, m.view(b))
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].BookingID.String() < all[j].BookingID.String() })
	start := min((f.Page-1)*f.Limit, len(all))
	end := min(start+f.Limit, len(all))
	return all[start:end], len(all), nil
}

func (m *MockBookingRepository) GetForUpdate(_ context.Context, id uuid.UUID, s domain.BookingScope) (*domain.LockedBooking, error) {
	b, ok := m.Bookings[id]
	if !ok || !bookingInScope(b, s) {
		return nil, domain.ErrBookingNotFound
	}
	return &domain.LockedBooking{
		BookingID: b.BookingID, TourID: b.TourID, AgencyID: b.AgencyID, UserID: b.UserID,
		NumberOfPeople: b.NumberOfPeople, Status: b.Status, TourStartDate: b.TourStartDate,
	}, nil
}

func (m *MockBookingRepository) UpdateStatus(_ context.Context, id uuid.UUID, status, reason string) error {
	b, ok := m.Bookings[id]
	if !ok {
		return domain.ErrBookingNotFound
	}
	b.Status = status
	b.CancellationReason = reason
	return nil
}

func (m *MockBookingRepository) CancelActiveForTour(_ context.Context, tourID uuid.UUID, reason string) ([]uuid.UUID, int, error) {
	ids := []uuid.UUID{}
	seats := 0
	for _, b := range m.Bookings {
		if b.TourID == tourID && (b.Status == domain.BookingPending || b.Status == domain.BookingConfirmed) {
			b.Status = domain.BookingCancelled
			b.CancellationReason = reason
			ids = append(ids, b.BookingID)
			seats += b.NumberOfPeople
		}
	}
	return ids, seats, nil
}

func (m *MockBookingRepository) GetOrCreateCustomerByUser(context.Context, uuid.UUID) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (m *MockBookingRepository) CreateCustomer(_ context.Context, c *domain.Customer) error {
	c.CustomerID = uuid.New()
	return nil
}

// MockPaymentRepository records payment amounts and status per booking.
type MockPaymentRepository struct {
	Amounts []int
	Status  map[uuid.UUID]string
}

func NewMockPaymentRepository() *MockPaymentRepository {
	return &MockPaymentRepository{Status: map[uuid.UUID]string{}}
}

func (m *MockPaymentRepository) Create(_ context.Context, p *domain.Payment) error {
	m.Amounts = append(m.Amounts, p.Amount)
	m.Status[p.BookingID] = "pending"
	return nil
}

func (m *MockPaymentRepository) SetStatusForBookings(_ context.Context, ids []uuid.UUID, from, to string) error {
	for _, id := range ids {
		if m.Status[id] == from {
			m.Status[id] = to
		}
	}
	return nil
}
