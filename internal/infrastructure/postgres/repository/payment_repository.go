package repository

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (p *paymentRepositoryDB) SetStatusForBookings(ctx context.Context, bookingIDs []uuid.UUID, from, to string) error {
	if len(bookingIDs) == 0 {
		return nil
	}
	ids := make([]string, len(bookingIDs))
	for i, id := range bookingIDs {
		ids[i] = id.String()
	}
	query := `UPDATE payments SET status = $1 WHERE booking_id = ANY($2::uuid[]) AND status = $3`
	_, err := p.executor(ctx).ExecContext(ctx, query, to, pq.Array(ids), from)
	return err
}

func (p *paymentRepositoryDB) executor(ctx context.Context) sqlx.ExtContext {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return p.db
}
