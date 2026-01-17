package repository

import (
	"database/sql"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TourRepositoryDB struct {
	db *sqlx.DB
}

func NewTourRepositoryDB(db *sqlx.DB) *TourRepositoryDB {
	return &TourRepositoryDB{
		db: db,
	}
}

func(h *TourRepositoryDB) CreateTour(tour *domain.Tour) error {
	query := `INSERT INTO tours (agencyID,name,startDate,endDate,description,lastEnrollmentDate,price,discount) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id;`

	return h.db.QueryRow(query,tour.AgencyID,tour.Name,tour.StartDate,tour.EndDate,tour.Description,tour.LastEnrollmentDate,tour.Price,tour.Discount).Scan(&tour.ID)
}

func(h *TourRepositoryDB)ListTour(agencyID int) ([]*domain.Tour,error) {
	var tours []*domain.Tour
	query := `SELECT name,startDate,endDate,description,lastEnrollmentDate,price,discount FROM tours WHERE agencyID=$1;`
	err := h.db.Select(&tours,query,agencyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil,nil
		}
		return nil,err
	}
	return tours,nil
}

func(h *TourRepositoryDB)UpdateTour(t domain.Tour) (*domain.Tour,error) {
	query := `UPDATE tours SET agencyID=$1,name=$2,startDate=$3,endDate=$4,description=$5,lastEnrollmentDate=$6,price=$7,discount=$8;`
	row := h.db.QueryRow(query,t.AgencyID,t.Name,t.StartDate,t.EndDate,t.Description,t.LastEnrollmentDate,t.Price,t.Discount)
	err := row.Err()
	if err != nil {
		return nil,err
	}
	return &t,nil
}

func(h *TourRepositoryDB)DeleteTour(tourID int) error {
	query := `DELETE FROM tours WHERE id=$1;`
	_,err := h.db.Exec(query,tourID)
	if err != nil {
		return err
	}
	return nil
}