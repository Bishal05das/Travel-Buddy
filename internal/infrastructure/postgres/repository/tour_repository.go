package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

// activeTourImageJoin attaches the tour's current image (if any) as image_path.
const activeTourImageJoin = `LEFT JOIN tour_images ti ON ti.tour_id = t.tour_id AND ti.is_active = TRUE AND ti.deleted_at IS NULL`

func (h *tourRepositoryDB) CreateTour(ctx context.Context, tour *domain.Tour) error {
	tx, err := h.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback() // no-op after a successful commit

	query := `INSERT INTO tours (agency_id,name,start_date,end_date,available_seat,description,last_enrollment_date,price,discount) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING tour_id;`
	err = tx.QueryRowContext(ctx, query, tour.AgencyID, tour.Name, tour.StartDate, tour.EndDate, tour.AvailableSeat, tour.Description, tour.LastEnrollmentDate, tour.Price, tour.Discount).Scan(&tour.TourID)
	if err != nil {
		return fmt.Errorf("insert tour: %w", err)
	}

	if tour.ImagePath != "" {
		imageQ := `INSERT INTO tour_images (tour_id, image_path, is_active) VALUES ($1, $2, TRUE);`
		if _, err := tx.ExecContext(ctx, imageQ, tour.TourID, tour.ImagePath); err != nil {
			return fmt.Errorf("insert tour image: %w", err)
		}
	}

	return tx.Commit()
}

func (h *tourRepositoryDB) ListTour(ctx context.Context, agencyID uuid.UUID, page, limit int) ([]*domain.Tour, error) {

	offset := (page - 1) * limit

	var tours []*domain.Tour
	query := `SELECT t.tour_id,t.agency_id,t.name,t.start_date,t.end_date,t.available_seat,t.description,t.last_enrollment_date,t.price,t.discount,t.status,COALESCE(ti.image_path,'') AS image_path FROM tours t ` + activeTourImageJoin + ` WHERE t.agency_id=$1 ORDER BY t.start_date DESC LIMIT $2 OFFSET $3;`
	err := h.db.SelectContext(ctx, &tours, query, agencyID, limit, offset)
	if err != nil {
		return nil, err
	}

	return tours, nil
}

func (h *tourRepositoryDB) Count(ctx context.Context, agencyID uuid.UUID) (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM tours WHERE agency_id=$1;`

	err := h.db.GetContext(ctx, &count, query, agencyID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (h *tourRepositoryDB) UpdateTour(ctx context.Context, t *domain.Tour) error {
	// agency_id is part of the filter, not the SET list: a tour can only be
	// updated through the agency that owns it and can never be moved.
	query := `UPDATE tours SET name=$1,start_date=$2,end_date=$3,available_seat=$4,description=$5,last_enrollment_date=$6,price=$7,discount=$8,updated_at=$9 WHERE tour_id=$10 AND agency_id=$11;`
	res, err := h.db.ExecContext(ctx, query, t.Name, t.StartDate, t.EndDate, t.AvailableSeat, t.Description, t.LastEnrollmentDate, t.Price, t.Discount, t.UpdatedAt, t.TourID, t.AgencyID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("tour not found")
	}
	return nil
}

func (h *tourRepositoryDB) DeleteTour(ctx context.Context, tourID uuid.UUID, agencyScope *uuid.UUID) error {
	query := `DELETE FROM tours WHERE tour_id=$1 AND ($2::uuid IS NULL OR agency_id=$2);`
	res, err := h.db.ExecContext(ctx, query, tourID, agencyScope)
	if err != nil {
		return err
	}
	return requireAffected(res, "tour not found")
}

func (h *tourRepositoryDB) GetByID(ctx context.Context, tourID uuid.UUID) (*domain.Tour, error) {
	query := `SELECT t.tour_id,t.agency_id,t.name,t.start_date,t.end_date,t.available_seat,t.description,t.last_enrollment_date,t.price,t.discount,t.status,COALESCE(ti.image_path,'') FROM tours t ` + activeTourImageJoin + ` WHERE t.tour_id=$1;`

	tour := &domain.Tour{}
	err := h.executor(ctx).QueryRowxContext(ctx, query, tourID).Scan(&tour.TourID, &tour.AgencyID, &tour.Name, &tour.StartDate, &tour.EndDate, &tour.AvailableSeat, &tour.Description, &tour.LastEnrollmentDate, &tour.Price, &tour.Discount, &tour.Status, &tour.ImagePath)
	if err == sql.ErrNoRows {
		return nil, errors.New("tour not found")
	}
	if err != nil {
		return nil, err
	}
	return tour, nil
}

func (h *tourRepositoryDB) UpdateAvailableSeats(ctx context.Context, tourID uuid.UUID, seats int) error {
	query := `UPDATE tours SET available_seat = $1, updated_at = CURRENT_TIMESTAMP WHERE tour_id = $2;`
	res, err := h.executor(ctx).ExecContext(ctx, query, seats, tourID)
	if err != nil {
		return err
	}
	return requireAffected(res, "tour not found")
}

func (h *tourRepositoryDB) GetByIDForUpdate(ctx context.Context, tourID uuid.UUID) (*domain.Tour, error) {
	// SELECT ... FOR UPDATE locks the row
	query := `SELECT tour_id,agency_id,name,start_date,end_date,available_seat,description,last_enrollment_date,price,discount,status FROM tours WHERE tour_id=$1 FOR UPDATE;`

	tour := &domain.Tour{}
	err := h.executor(ctx).QueryRowxContext(ctx, query, tourID).Scan(&tour.TourID, &tour.AgencyID, &tour.Name, &tour.StartDate, &tour.EndDate, &tour.AvailableSeat, &tour.Description, &tour.LastEnrollmentDate, &tour.Price, &tour.Discount, &tour.Status)
	if err == sql.ErrNoRows {
		return nil, errors.New("tour not found")
	}
	if err != nil {
		return nil, err
	}
	return tour, nil
}

func (h *tourRepositoryDB) UpdateTourStatus(ctx context.Context, tourID uuid.UUID, status string, agencyScope *uuid.UUID) error {
	query := `UPDATE tours SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE tour_id= $2 AND ($3::uuid IS NULL OR agency_id=$3);`
	res, err := h.executor(ctx).ExecContext(ctx, query, status, tourID, agencyScope)
	if err != nil {
		return err
	}
	return requireAffected(res, "tour not found")
}

func (h *tourRepositoryDB) executor(ctx context.Context) sqlx.ExtContext {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return h.db
}
