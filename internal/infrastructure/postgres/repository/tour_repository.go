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

type tourRepositoryDB struct {
	db *sqlx.DB
}

func NewTourRepositoryDB(db *sqlx.DB) port.TourRepository {
	return &tourRepositoryDB{
		db: db,
	}
}

func (h *tourRepositoryDB) CreateTour(ctx context.Context, tour *domain.Tour) error {
	query := `INSERT INTO tours (agency_id,name,start_date,end_date,available_seat,description,last_enrollment_date,price,discount) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING tour_id;`

	return h.db.QueryRowContext(ctx, query, tour.AgencyID, tour.Name, tour.StartDate, tour.EndDate, tour.AvailableSeat, tour.Description, tour.LastEnrollmentDate, tour.Price, tour.Discount).Scan(&tour.TourID)
}

func (h *tourRepositoryDB) ListTour(ctx context.Context, agencyID uuid.UUID) ([]*domain.Tour, error) {

	var tours []*domain.Tour
	query := `SELECT name,start_date,end_date,available_seat,description,last_enrollment_date,price,discount FROM tours WHERE agency_id=$1;`
	err := h.db.SelectContext(ctx, &tours, query, agencyID)
	if err != nil {
		return nil, err
	}

	return tours, nil
}

func (h *tourRepositoryDB) UpdateTour(ctx context.Context, t *domain.Tour) error {
	query := `UPDATE tours SET agency_id=$1,name=$2,start_date=$3,end_date=$4,available_seat=$5,description=$6,last_enrollment_date=$7,price=$8,discount=$9 WHERE tour_id=$10;`
	res,err := h.db.ExecContext(ctx, query, t.AgencyID, t.Name, t.StartDate, t.EndDate, t.AvailableSeat, t.Description, t.LastEnrollmentDate, t.Price, t.Discount, t.TourID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("tour not found")
	}
	return nil
}

func (h *tourRepositoryDB) DeleteTour(ctx context.Context, tourID uuid.UUID) error {
	query := `DELETE FROM tours WHERE tour_id=$1;`
	_, err := h.db.ExecContext(ctx, query, tourID)
	if err != nil {
		return err
	}
	return nil
}

func (h *tourRepositoryDB) GetByID(ctx context.Context, tourID uuid.UUID) (*domain.Tour, error) {
	query := `SELECT agency_id,name,start_date,end_date,available_seat,description,last_enrollment_date,price,discount,status FROM tours WHERE tour_id=$1;`

	tour := &domain.Tour{}
	err := h.executor(ctx).QueryRowxContext(ctx, query, tourID).Scan(&tour.AgencyID, &tour.Name, &tour.StartDate, &tour.EndDate, &tour.AvailableSeat, &tour.Description, &tour.LastEnrollmentDate, &tour.Price, &tour.Discount, &tour.Status)
	if err == sql.ErrNoRows {
		return nil, errors.New("tour not found")
	}
	if err != nil {
		return nil, err
	}
	return tour, nil
}

func (h *tourRepositoryDB) UpdateAvailableSeats(ctx context.Context, tourID uuid.UUID, seats int) error {
	query := `UPDATE tours SET available_seat = $1, updated_at = CURRENT_TIMESTAMP WHERE tour_id = $2`
	_, err := h.executor(ctx).ExecContext(ctx, query, seats, tourID)
	return err
}

func (h *tourRepositoryDB) GetByIDForUpdate(ctx context.Context, tourID uuid.UUID) (*domain.Tour, error) {
	// SELECT ... FOR UPDATE locks the row
	query := `SELECT agency_id,name,start_date,end_date,available_seat,description,last_enrollment_date,price,discount,status FROM tours WHERE tour_id=$1 FOR UPDATE;`

	tour := &domain.Tour{}
	err := h.executor(ctx).QueryRowxContext(ctx, query, tourID).Scan(&tour.AgencyID, &tour.Name, &tour.StartDate, &tour.EndDate, &tour.AvailableSeat, &tour.Description, &tour.LastEnrollmentDate, &tour.Price, &tour.Discount, &tour.Status)
	if err == sql.ErrNoRows {
		return nil, errors.New("tour not found")
	}
	if err != nil {
		return nil, err
	}
	return tour, nil
}

func (h *tourRepositoryDB) executor(ctx context.Context) sqlx.ExtContext {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return h.db
}
