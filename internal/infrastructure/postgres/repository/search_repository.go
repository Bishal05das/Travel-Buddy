package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/jmoiron/sqlx"
)

type searchRepository struct{ db *sqlx.DB }

func NewSearchRepository(db *sqlx.DB) port.SearchRepository { return &searchRepository{db: db} }

var tourSearchColumns = []searchColumn{{"t.name", 3}, {"a.name", 2}, {"t.description", 1}}

const searchTourFrom = ` FROM tours t JOIN agency a ON t.agency_id=a.agency_id `

// Cast before multiplication to preserve the booking formula without integer overflow.
const searchablePrice = `(t.price::bigint - (t.price::bigint * COALESCE(t.discount, 0)) / 100)`

func tourSearchConditions(filter domain.TourSearchFilter, sql *searchSQL) string {
	conditions := []string{"a.is_active=TRUE"}
	conditions = append(conditions, sql.matches(filter.Query, tourSearchColumns)...)
	if filter.MinPrice != nil {
		conditions = append(conditions, searchablePrice+">="+sql.bind(*filter.MinPrice)+"::numeric")
	}
	if filter.MaxPrice != nil {
		conditions = append(conditions, searchablePrice+"<="+sql.bind(*filter.MaxPrice)+"::numeric")
	}
	if filter.StartDate != nil {
		conditions = append(conditions, "t.start_date>="+sql.bind(*filter.StartDate))
	}
	if filter.EndDate != nil {
		conditions = append(conditions, "t.end_date<="+sql.bind(*filter.EndDate))
	}
	if filter.Status != nil {
		conditions = append(conditions, "t.status="+sql.bind(*filter.Status))
	}
	if filter.AgencyID != nil {
		conditions = append(conditions, "t.agency_id="+sql.bind(*filter.AgencyID))
	}
	return strings.Join(conditions, " AND ")
}

func (r *searchRepository) SearchTours(ctx context.Context, filter domain.TourSearchFilter) ([]domain.TourSearchResponse, error) {
	filter.Query = domain.NormalizeSearchQuery(filter.Query)
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	sql := &searchSQL{}
	where := tourSearchConditions(filter, sql)
	order := sql.ranking(filter.Query, tourSearchColumns)
	query := `SELECT t.tour_id,t.name,t.price,COALESCE(t.discount,0),t.start_date,t.end_date,
		t.status,t.available_seat,a.agency_id,a.name,t.description,t.last_enrollment_date,
		COALESCE(ti.image_path,'')` + searchTourFrom + activeTourImageJoin + " WHERE " + where +
		" ORDER BY " + order + "t.created_at DESC,t.tour_id ASC LIMIT " + sql.bind(filter.Limit) + " OFFSET " + sql.bind(filter.Offset)
	tours := []domain.TourSearchResponse{}
	err := r.withSearch(ctx, filter.Query, func(db sqlx.ExtContext) error {
		rows, err := db.QueryxContext(ctx, query, sql.args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var tour domain.TourSearchResponse
			if err := rows.Scan(&tour.TourID, &tour.Name, &tour.Price, &tour.Discount, &tour.StartDate, &tour.EndDate,
				&tour.Status, &tour.AvailableSeat, &tour.AgencyID, &tour.AgencyName, &tour.Description, &tour.LastEnrollmentDate, &tour.ImagePath); err != nil {
				return err
			}
			tours = append(tours, tour)
		}
		return rows.Err()
	})
	return tours, err
}

func (r *searchRepository) CountTours(ctx context.Context, filter domain.TourSearchFilter) (int, error) {
	filter.Query = domain.NormalizeSearchQuery(filter.Query)
	sql := &searchSQL{}
	query := "SELECT COUNT(*)" + searchTourFrom + " WHERE " + tourSearchConditions(filter, sql)
	var count int
	err := r.withSearch(ctx, filter.Query, func(db sqlx.ExtContext) error {
		return db.QueryRowxContext(ctx, query, sql.args...).Scan(&count)
	})
	return count, err
}

func (r *searchRepository) SearchAgencies(ctx context.Context, query string, limit int) ([]domain.Agency, error) {
	query = domain.NormalizeSearchQuery(query)
	agencies := []domain.Agency{}
	if query == "" {
		return agencies, nil
	}
	sql := &searchSQL{}
	columns := []searchColumn{{"a.name", 3}}
	where := strings.Join(append([]string{"a.is_active=TRUE"}, sql.matches(query, columns)...), " AND ")
	order := sql.ranking(query, columns)
	statement := `SELECT a.agency_id,a.name,COALESCE(a.address,''),COALESCE(a.reg_id,''),
		COALESCE(a.rating,0),a.is_active,a.created_at,a.updated_at,COALESCE(ai.image_path,'')
		FROM agency a LEFT JOIN agency_images ai ON ai.agency_id=a.agency_id AND ai.is_active=TRUE AND ai.deleted_at IS NULL
		WHERE ` + where + " ORDER BY " + order + "a.created_at DESC,a.agency_id ASC LIMIT " + sql.bind(limit)
	err := r.withSearch(ctx, query, func(db sqlx.ExtContext) error {
		rows, err := db.QueryxContext(ctx, statement, sql.args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var agency domain.Agency
			if err := rows.Scan(&agency.AgencyID, &agency.Name, &agency.Address, &agency.RegistrationID, &agency.Rating,
				&agency.IsActive, &agency.CreatedAt, &agency.UpdatedAt, &agency.ImagePath); err != nil {
				return err
			}
			agencies = append(agencies, agency)
		}
		return rows.Err()
	})
	return agencies, err
}

type searchSQL struct{ args []any }

func (s *searchSQL) bind(value any) string {
	s.args = append(s.args, value)
	return fmt.Sprintf("$%d", len(s.args))
}
