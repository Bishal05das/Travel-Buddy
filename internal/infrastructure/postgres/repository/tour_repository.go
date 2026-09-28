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

	// A new tour has no bookings, so every seat is available.
	query := `INSERT INTO tours (agency_id,name,start_date,end_date,total_seat,available_seat,description,last_enrollment_date,price,discount) VALUES ($1,$2,$3,$4,$5,$5,$6,$7,$8,$9) RETURNING tour_id, available_seat;`
	err = tx.QueryRowContext(ctx, query, tour.AgencyID, tour.Name, tour.StartDate, tour.EndDate, tour.TotalSeat, tour.Description, tour.LastEnrollmentDate, tour.Price, tour.Discount).Scan(&tour.TourID, &tour.AvailableSeat)
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

	tours := []*domain.Tour{} // encode an empty page as [], not null
	query := `SELECT t.tour_id,t.agency_id,t.name,t.start_date,t.end_date,t.total_seat,t.available_seat,t.description,t.last_enrollment_date,t.price,t.discount,t.status,COALESCE(ti.image_path,'') AS image_path FROM tours t ` + activeTourImageJoin + ` WHERE t.agency_id=$1 ORDER BY t.start_date DESC LIMIT $2 OFFSET $3;`
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
	//
	// The client sets the capacity (total_seat). The remaining seats are
	// recomputed from the seats already booked (total_seat - available_seat,
	// read from the row being updated), so an edit never discards bookings.
	// The capacity guard is in the WHERE clause, so the check and the write are
	// one atomic statement under the same row lock bookings take.
	query := `
	UPDATE tours SET
		name=$1, start_date=$2, end_date=$3,
		available_seat = $4 - (total_seat - available_seat),
		total_seat = $4,
		description=$5, last_enrollment_date=$6, price=$7, discount=$8, updated_at=$9
	WHERE tour_id=$10 AND agency_id=$11 AND $4 >= total_seat - available_seat
	RETURNING available_seat;`
	err := h.db.QueryRowxContext(ctx, query, t.Name, t.StartDate, t.EndDate, t.TotalSeat, t.Description, t.LastEnrollmentDate, t.Price, t.Discount, t.UpdatedAt, t.TourID, t.AgencyID).Scan(&t.AvailableSeat)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// Nothing matched: either the tour does not exist for this agency, or the
	// new capacity is below the booked seats.
	var booked int
	err = h.db.QueryRowxContext(ctx, `SELECT total_seat - available_seat FROM tours WHERE tour_id=$1 AND agency_id=$2;`, t.TourID, t.AgencyID).Scan(&booked)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("tour not found")
	}
	if err != nil {
		return err
	}
	return fmt.Errorf("%w (%d booked)", domain.ErrCapacityBelowBooked, booked)
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
	query := `SELECT t.tour_id,t.agency_id,t.name,t.start_date,t.end_date,t.total_seat,t.available_seat,t.description,t.last_enrollment_date,t.price,t.discount,t.status,COALESCE(ti.image_path,'') FROM tours t ` + activeTourImageJoin + ` WHERE t.tour_id=$1;`

	tour := &domain.Tour{}
	err := h.executor(ctx).QueryRowxContext(ctx, query, tourID).Scan(&tour.TourID, &tour.AgencyID, &tour.Name, &tour.StartDate, &tour.EndDate, &tour.TotalSeat, &tour.AvailableSeat, &tour.Description, &tour.LastEnrollmentDate, &tour.Price, &tour.Discount, &tour.Status, &tour.ImagePath)
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
	query := `SELECT tour_id,agency_id,name,start_date,end_date,total_seat,available_seat,description,last_enrollment_date,price,discount,status FROM tours WHERE tour_id=$1 FOR UPDATE;`

	tour := &domain.Tour{}
	err := h.executor(ctx).QueryRowxContext(ctx, query, tourID).Scan(&tour.TourID, &tour.AgencyID, &tour.Name, &tour.StartDate, &tour.EndDate, &tour.TotalSeat, &tour.AvailableSeat, &tour.Description, &tour.LastEnrollmentDate, &tour.Price, &tour.Discount, &tour.Status)
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
