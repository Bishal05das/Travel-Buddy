package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/jmoiron/sqlx"
)

type searchRepository struct {
	db *sqlx.DB
}

func NewSearchRepository(db *sqlx.DB) port.SearchRepository {
	return &searchRepository{db: db}
}

func (r *searchRepository) SearchTours(ctx context.Context,filter domain.TourSearchFilter,) ([]domain.TourSearchResponse, error) {

	baseQuery := `
	SELECT 
		t.tour_id,
		t.name,
		t.price,
		t.start_date,
		t.end_date,
		t.status,
		t.available_seat,
		a.agency_id,
		a.name AS agency_name
	FROM tours t
	JOIN agency a ON t.agency_id = a.agency_id
	WHERE 1=1
	`

	var args []interface{}
	var conditions []string
	argPos := 1

	// Search text
	if filter.Query != "" {
		conditions = append(conditions,
			fmt.Sprintf(`(
				t.name ILIKE $%d OR
				t.description ILIKE $%d OR
				a.name ILIKE $%d
			)`, argPos, argPos, argPos),
		)
		args = append(args, "%"+filter.Query+"%")
		argPos++
	}

	if filter.MinPrice != nil {
		conditions = append(conditions,
			fmt.Sprintf("t.price >= $%d", argPos))
		args = append(args, *filter.MinPrice)
		argPos++
	}

	if filter.MaxPrice != nil {
		conditions = append(conditions,
			fmt.Sprintf("t.price <= $%d", argPos))
		args = append(args, *filter.MaxPrice)
		argPos++
	}

	if filter.StartDate != nil {
		conditions = append(conditions,
			fmt.Sprintf("t.start_date >= $%d", argPos))
		args = append(args, *filter.StartDate)
		argPos++
	}

	if filter.EndDate != nil {
		conditions = append(conditions,
			fmt.Sprintf("t.end_date <= $%d", argPos))
		args = append(args, *filter.EndDate)
		argPos++
	}

	// if filter.Status != nil {
	// 	conditions = append(conditions,
	// 		fmt.Sprintf("t.status = $%d", argPos))
	// 	args = append(args, *filter.Status)
	// 	argPos++
	// }

	// if filter.AgencyID != nil {
	// 	conditions = append(conditions,
	// 		fmt.Sprintf("t.agency_id = $%d", argPos))
	// 	args = append(args, *filter.AgencyID)
	// 	argPos++
	// }

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	baseQuery += fmt.Sprintf(
		" ORDER BY t.created_at DESC LIMIT $%d OFFSET $%d",
		argPos, argPos+1,
	)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tours []domain.TourSearchResponse

	for rows.Next() {
		var t domain.TourSearchResponse
		err := rows.Scan(
			&t.TourID,
			&t.Name,
			&t.Price,
			&t.StartDate,
			&t.EndDate,
			&t.Status,
			&t.AvailableSeat,
			&t.AgencyID,
			&t.AgencyName,
		)
		if err != nil {
			return nil, err
		}
		tours = append(tours, t)
	}

	return tours, nil
}

func (r *searchRepository) SearchAgencies(
	ctx context.Context,
	query string,
	limit int,
) ([]domain.Agency, error) {

	sqlQuery := `
	SELECT agency_id, name, rating
	FROM agency
	WHERE name ILIKE $1
	ORDER BY created_at DESC
	LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, sqlQuery, "%"+query+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agencies []domain.Agency

	for rows.Next() {
		var a domain.Agency
		if err := rows.Scan(&a.AgencyID, &a.Name, &a.Rating); err != nil {
			return nil, err
		}
		agencies = append(agencies, a)
	}

	return agencies, nil
}