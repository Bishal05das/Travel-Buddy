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
	query := `INSERT INTO tours (agency_id,name,start_date,end_date,available_seat,description,last_enrollment_date,price,discount) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id;`

	return h.db.QueryRow(query,tour.AgencyID,tour.Name,tour.StartDate,tour.EndDate,tour.AvailableSeat,tour.Description,tour.LastEnrollmentDate,tour.Price,tour.Discount).Scan(&tour.ID)
}

func(h *TourRepositoryDB)ListTour(agencyID int) ([]*domain.Tour,error) {
	var tours []*domain.Tour
	query := `SELECT name,start_date,end_date,available_seat,description,last_enrollment_date,price,discount FROM tours WHERE agency_id=$1;`
	err := h.db.Select(&tours,query,agencyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil,nil
		}
		return nil,err
	}
	return tours,nil
}

func(h *TourRepositoryDB)UpdateTour(t *domain.Tour) error {
	query := `UPDATE tours SET agency_id=$1,name=$2,start_date=$3,end_date=$4,available_seat=$5,description=$6,last_enrollment_date=$7,price=$8,discount=$9;`
	row := h.db.QueryRow(query,t.AgencyID,t.Name,t.StartDate,t.EndDate,t.AvailableSeat,t.Description,t.LastEnrollmentDate,t.Price,t.Discount)
	err := row.Err()
	if err != nil {
		return err
	}
	return nil
}

func(h *TourRepositoryDB)DeleteTour(tourID int) error {
	query := `DELETE FROM tours WHERE id=$1;`
	_,err := h.db.Exec(query,tourID)
	if err != nil {
		return err
	}
	return nil
}