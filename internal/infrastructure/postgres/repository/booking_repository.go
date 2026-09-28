package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type bookingRepository struct {
	db *sqlx.DB
}

func NewBookingRepository(db *sqlx.DB) port.BookingRepository {
	return &bookingRepository{
		db: db,
	}
}

func (r *bookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	query := `
		INSERT INTO bookings (customer_id, user_id, member_id, tour_id, number_of_people, total_price, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING booking_id, booking_date`

	err := r.executor(ctx).QueryRowxContext(
		ctx, query,
		booking.CustomerID, booking.UserID, booking.MemberID, booking.TourID,
		booking.NumberOfPeople, booking.TotalPrice, booking.Status,
	).Scan(&booking.BookingID, &booking.BookingDate)

	return err
}

// bookingSelect returns bookings with customer, tour, agency and payment
// details. Every filter is optional (NULL / ” means "any"):
// $1 booking_id, $2 agency_id, $3 user_id, $4 tour_id, $5 status.
const bookingSelect = `
	SELECT
		b.booking_id, b.status, b.number_of_people, b.total_price, b.booking_date,
		CASE WHEN b.user_id IS NOT NULL THEN 'user' ELSE 'agency_member' END,
		b.user_id, b.member_id,
		b.customer_id,
		COALESCE(c.name, u.name, ''), COALESCE(c.email, u.email, ''), COALESCE(c.phone, u.phone, ''),
		t.tour_id, t.name, t.start_date, a.agency_id, a.name,
		COALESCE(p.method, ''), COALESCE(p.transaction_id, ''), COALESCE(p.status, ''),
		b.created_at, b.updated_at,
		COUNT(*) OVER ()
	FROM bookings b
	JOIN tours t ON t.tour_id = b.tour_id
	JOIN agency a ON a.agency_id = t.agency_id
	JOIN customers c ON c.customer_id = b.customer_id
	LEFT JOIN users u ON u.user_id = c.user_id
	LEFT JOIN LATERAL (
		SELECT method, transaction_id, status
		FROM payments
		WHERE booking_id = b.booking_id
		ORDER BY created_at DESC
		LIMIT 1
	) p ON TRUE
	WHERE ($1::uuid IS NULL OR b.booking_id = $1)
	  AND ($2::uuid IS NULL OR t.agency_id = $2)
	  AND ($3::uuid IS NULL OR b.user_id = $3)
	  AND ($4::uuid IS NULL OR b.tour_id = $4)
	  AND ($5::text = '' OR b.status = $5)
	ORDER BY b.created_at DESC
	LIMIT $6 OFFSET $7`

func (r *bookingRepository) queryBookings(ctx context.Context, bookingID *uuid.UUID, f domain.BookingFilter, limit, offset int) ([]*domain.BookingResponse, int, error) {
	rows, err := r.executor(ctx).QueryxContext(ctx, bookingSelect,
		bookingID, f.AgencyID, f.UserID, f.TourID, f.Status, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	bookings := []*domain.BookingResponse{}
	total := 0
	for rows.Next() {
		b := &domain.BookingResponse{}
		if err := rows.Scan(
			&b.BookingID, &b.Status, &b.NumberOfPeople, &b.TotalPrice, &b.BookingDate,
			&b.CreatedBy, &b.UserID, &b.MemberID,
			&b.CustomerID, &b.CustomerName, &b.CustomerEmail, &b.CustomerPhone,
			&b.TourID, &b.TourName, &b.TourStartDate, &b.AgencyID, &b.AgencyName,
			&b.PaymentMethod, &b.TransactionID, &b.PaymentStatus,
			&b.CreatedAt, &b.UpdatedAt,
			&total,
		); err != nil {
			return nil, 0, err
		}
		bookings = append(bookings, b)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return bookings, total, nil
}

func (r *bookingRepository) GetByID(ctx context.Context, id uuid.UUID, scope domain.BookingScope) (*domain.BookingResponse, error) {
	bookings, _, err := r.queryBookings(ctx, &id, domain.BookingFilter{BookingScope: scope}, 1, 0)
	if err != nil {
		return nil, err
	}
	if len(bookings) == 0 {
		return nil, domain.ErrBookingNotFound
	}
	return bookings[0], nil
}

func (r *bookingRepository) List(ctx context.Context, f domain.BookingFilter) ([]*domain.BookingResponse, int, error) {
	bookings, total, err := r.queryBookings(ctx, nil, f, f.Limit, (f.Page-1)*f.Limit)
	if err != nil {
		return nil, 0, err
	}
	if len(bookings) == 0 && f.Page > 1 {
		// A page past the end has no rows to carry the window count.
		var count int
		q := `SELECT COUNT(*) FROM bookings b JOIN tours t ON t.tour_id = b.tour_id
		      WHERE ($1::uuid IS NULL OR t.agency_id = $1) AND ($2::uuid IS NULL OR b.user_id = $2)
		        AND ($3::uuid IS NULL OR b.tour_id = $3) AND ($4::text = '' OR b.status = $4)`
		if err := r.executor(ctx).QueryRowxContext(ctx, q, f.AgencyID, f.UserID, f.TourID, f.Status).Scan(&count); err != nil {
			return nil, 0, err
		}
		total = count
	}
	return bookings, total, nil
}

func (r *bookingRepository) GetForUpdate(ctx context.Context, id uuid.UUID, scope domain.BookingScope) (*domain.LockedBooking, error) {
	query := `
		SELECT b.booking_id, b.tour_id, t.agency_id, b.user_id, b.number_of_people, b.status, t.start_date
		FROM bookings b
		JOIN tours t ON t.tour_id = b.tour_id
		WHERE b.booking_id = $1
		  AND ($2::uuid IS NULL OR t.agency_id = $2)
		  AND ($3::uuid IS NULL OR b.user_id = $3)
		FOR UPDATE OF b`
	lb := &domain.LockedBooking{}
	err := r.executor(ctx).QueryRowxContext(ctx, query, id, scope.AgencyID, scope.UserID).Scan(
		&lb.BookingID, &lb.TourID, &lb.AgencyID, &lb.UserID, &lb.NumberOfPeople, &lb.Status, &lb.TourStartDate)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrBookingNotFound
	}
	if err != nil {
		return nil, err
	}
	return lb, nil
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE bookings SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE booking_id = $2`
	res, err := r.executor(ctx).ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}
	if err := requireAffected(res, "booking not found"); err != nil {
		return domain.ErrBookingNotFound
	}
	return nil
}

func (r *bookingRepository) GetOrCreateCustomerByUser(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	var customerID uuid.UUID

	query := `SELECT customer_id FROM customers WHERE user_id = $1`
	err := r.executor(ctx).QueryRowxContext(ctx, query, userID).Scan(&customerID)
	if err == nil {
		return customerID, nil
	}

	if err != sql.ErrNoRows {
		return uuid.Nil, err
	}

	// Create new customer
	insertQuery := `
		INSERT INTO customers (user_id)
		VALUES ($1)
		RETURNING customer_id`

	err = r.executor(ctx).QueryRowxContext(ctx, insertQuery, userID).Scan(&customerID)
	if err != nil {
		return uuid.Nil, err
	}
	return customerID, nil
}

func (r *bookingRepository) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	query := `
		INSERT INTO customers (user_id, name, email, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING customer_id`

	return r.executor(ctx).QueryRowxContext(
		ctx, query,
		customer.UserID, customer.Name, customer.Email, customer.Phone,
	).Scan(&customer.CustomerID)
}

func (r *bookingRepository) executor(ctx context.Context) sqlx.ExtContext {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return r.db
}
