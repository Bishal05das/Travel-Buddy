package repository

import (
	"database/sql"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/jmoiron/sqlx"
)

type agencyRepositoryDB struct {
	db *sqlx.DB
}

func NewAgencyRepositoryDB(db *sqlx.DB) port.AgencyRepository {
	return &agencyRepositoryDB{
		db: db,
	}
}

func(h *agencyRepositoryDB) CreateAgency(Agency *domain.Agency) error {
	query := `INSERT INTO agency (name,address,registrationid) VALUES ($1,$2,$3) RETURNING id;`

	return h.db.QueryRow(query,Agency.Name,Agency.Address,Agency.RegistrationID).Scan(&Agency.AgencyID)
}

func(h *agencyRepositoryDB)ListAgency(agencyID int) ([]*domain.Agency,error) {
	var agencys []*domain.Agency
	query := `SELECT name,rating FROM Agency WHERE agencyID=$1;`
	err := h.db.Select(&agencys,query,agencyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil,nil
		}
		return nil,err
	}
	return agencys,nil
}

func(h *agencyRepositoryDB)UpdateAgency(agency *domain.Agency) error {
	query := `UPDATE Agencys SET name=$1,address=$2,registration_id=$3;`
	row := h.db.QueryRow(query,agency.Name,agency.Address,agency.RegistrationID)
	err := row.Err()
	if err != nil {
		return err
	}
	return nil
}

func(h *agencyRepositoryDB)DeleteAgency(agencyID int) error {
	query := `DELETE FROM Agencys WHERE id=$1;`
	_,err := h.db.Exec(query,agencyID)
	if err != nil {
		return err
	}
	return nil
}