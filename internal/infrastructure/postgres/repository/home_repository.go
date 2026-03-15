package repository

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/jmoiron/sqlx"
)

type homeRepositoryDB struct {
	db *sqlx.DB
}

func NewHomeRepositoryDB(db *sqlx.DB) port.HomeRepository {
	return &homeRepositoryDB{
		db: db,
	}
}

func (r *homeRepositoryDB) GetTopTours(ctx context.Context, limit int) ([]domain.HomeTour, error) {
	query := `
		SELECT
			t.tour_id,
			t.name,
			a.agency_id,
			a.name,
			a.rating,
			t.start_date,
			t.end_date,
			t.price,
			t.discount,
			(t.price - t.discount)       AS final_price,
			t.available_seat,
			t.last_enrollment_date,
			t.description,
			t.status,
			COUNT(b.booking_id)          AS total_bookings
		FROM tours t
		JOIN agency a ON t.agency_id = a.agency_id
		LEFT JOIN bookings b
			ON t.tour_id = b.tour_id
			AND b.status NOT IN ('cancelled')
		WHERE t.status = 'open'
		  AND t.available_seat > 0
		  AND t.last_enrollment_date >= CURRENT_DATE
		GROUP BY t.tour_id, a.agency_id, a.name, a.rating
		ORDER BY total_bookings DESC, t.start_date ASC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tours []domain.HomeTour
	for rows.Next() {
		var t domain.HomeTour
		if err := rows.Scan(
			&t.TourID, &t.Name, &t.AgencyID,
			&t.AgencyName, &t.AgencyRating,
			&t.StartDate, &t.EndDate,
			&t.Price, &t.Discount, &t.FinalPrice,
			&t.AvailableSeat, &t.LastEnrollmentDate,
			&t.Description, &t.Status,
			&t.TotalBookings,
		); err != nil {
			return nil, err
		}
		tours = append(tours, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if tours == nil {
		tours = []domain.HomeTour{}
	}
	return tours, nil
}

func (r *homeRepositoryDB) GetTopAgencies(ctx context.Context, limit int) ([]domain.HomeAgency, error) {
	query := `
		SELECT
			a.agency_id,
			a.name,
			COALESCE(a.address, ''),
			a.rating,
			COUNT(t.tour_id) AS total_tours
		FROM agency a
		LEFT JOIN tours t
			ON a.agency_id = t.agency_id
			AND t.status = 'open'
		WHERE a.is_active = true
		GROUP BY a.agency_id, a.name, a.address, a.rating
		ORDER BY a.rating DESC, total_tours DESC
		LIMIT $1
	`
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agencies []domain.HomeAgency
	for rows.Next() {
		var ag domain.HomeAgency
		if err := rows.Scan(
			&ag.AgencyID, &ag.Name, &ag.Address,
			&ag.Rating, &ag.TotalTours,
		); err != nil {
			return nil, err
		}
		agencies = append(agencies, ag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if agencies == nil {
		agencies = []domain.HomeAgency{}
	}
	return agencies, nil
}
