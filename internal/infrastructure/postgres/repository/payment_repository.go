package repository

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type paymentRepositoryDB struct {
	db *sqlx.DB
}

func NewPaymentRepositoryDB(db *sqlx.DB) port.PaymentRepository {
	return &paymentRepositoryDB{
		db: db,
	}
}

func (p *paymentRepositoryDB) Create(ctx context.Context, payment *domain.Payment) error {
	query := `INSERT INTO payments (booking_id,transaction_id,amount,method) VALUES ($1,$2,$3,$4) RETURNING payment_id;`

	return p.executor(ctx).QueryRowxContext(ctx, query, payment.BookingID, payment.TransactionID, payment.Amount, payment.Method).Scan(&payment.PaymentID)
}

func (p *paymentRepositoryDB) SetStatusForBooking(ctx context.Context, bookingID uuid.UUID, from, to string) error {
	query := `UPDATE payments SET status = $1 WHERE booking_id = $2 AND status = $3`
	_, err := p.executor(ctx).ExecContext(ctx, query, to, bookingID, from)
	return err
}

func (p *paymentRepositoryDB) executor(ctx context.Context) sqlx.ExtContext {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return p.db
}
