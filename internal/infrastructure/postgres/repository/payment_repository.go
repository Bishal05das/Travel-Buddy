package repository

import (
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
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

func (p *paymentRepositoryDB) Create(payment *domain.Payment) error{
	query := `INSERT INTO payments (booking_id,transaction_id,amount,method) VALUES ($1,$2,$3,$4) RETURNING payment_id;`

	return p.db.QueryRow(query,payment.BookingID,payment.TransactionID,payment.Amount,payment.Method).Scan(&payment.PaymentID)
}