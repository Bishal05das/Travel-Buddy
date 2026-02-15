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

	err := r.db.QueryRowContext(
		ctx, query,
		booking.CustomerID, booking.UserID, booking.MemberID, booking.TourID,
		booking.NumberOfPeople, booking.TotalPrice, booking.Status,
	).Scan(&booking.BookingID, &booking.BookingDate)

	return err
}

func (r *bookingRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.BookingResponse, error) {
	query := `
		SELECT b.booking_id, b.customer_id, 
			COALESCE(c.name, u.name) as customer_name,
			b.tour_id, t.name as tour_name, a.name as agency_name,
			b.booking_date, b.number_of_people, b.total_price, b.status,
			CASE 
				WHEN b.user_id IS NOT NULL THEN 'user'
				WHEN b.member_id IS NOT NULL THEN 'agency_member'
			END as created_by,
			b.created_at
		FROM bookings b
		JOIN tours t ON b.tour_id = t.tour_id
		JOIN agency a ON t.agency_id = a.agency_id
		JOIN customers c ON b.customer_id = c.customer_id
		LEFT JOIN users u ON c.user_id = u.user_id
		WHERE b.booking_id = $1`

	booking := &domain.BookingResponse{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.BookingID, &booking.CustomerID, &booking.CustomerName,
		&booking.TourID, &booking.TourName, &booking.AgencyName,
		&booking.BookingDate, &booking.NumberOfPeople, &booking.TotalPrice,
		&booking.Status, &booking.CreatedBy, &booking.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("booking not found")
	}
	if err != nil {
		return nil, err
	}

	return booking, nil
}
func (r *bookingRepository) Update(ctx context.Context, booking *domain.Booking) error {
	query := `
		UPDATE bookings
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE booking_id = $2
		RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query, booking.Status, booking.BookingID).Scan(&booking.UpdatedAt)

	if err == sql.ErrNoRows {
		return errors.New("booking not found")
	}
	return err
}

func (r *bookingRepository) Cancel(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE bookings SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP WHERE booking_id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("booking not found")
	}

	return nil
}

func (r *bookingRepository) GetOrCreateCustomerByUser(ctx context.Context, userID uuid.UUID) (uuid.UUID,error) {
	var customerID uuid.UUID

	query := `SELECT customer_id FROM customers WHERE user_id = $1`
	err := r.db.QueryRowContext(ctx,query,userID).Scan(&customerID)
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

	err = r.db.QueryRowContext(ctx, insertQuery, userID).Scan(&customerID)
	return customerID, err
}

func (r *bookingRepository) CreateCustomer(ctx context.Context, customer *domain.Customer) error {
	query := `
		INSERT INTO customers (user_id, name, email, phone)
		VALUES ($1, $2, $3, $4)
		RETURNING customer_id`

	return r.db.QueryRowContext(
		ctx, query,
		customer.UserID, customer.Name, customer.Email, customer.Phone,
	).Scan(&customer.CustomerID)
}