package repository

import (
	"context"
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
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

func (h *agencyRepositoryDB) CreateAgency(ctx context.Context, Agency *domain.Agency) error {
	query := `INSERT INTO agency (name,address,reg_id) VALUES ($1,$2,$3) RETURNING agency_id;`

	return h.db.QueryRowContext(ctx, query, Agency.Name, Agency.Address, Agency.RegistrationID).Scan(&Agency.AgencyID)
}

// func (h *agencyRepositoryDB) ListAgency(ctx context.Context,agencyID int) ([]*domain.Agency,error) {
// 	var agencys []*domain.Agency
// 	query := `SELECT name,rating FROM Agency WHERE agencyID=$1;`
// 	err := h.db.SelectContext(ctx, &agencys, query, agencyID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return agencys, nil
// }

func (h *agencyRepositoryDB) UpdateAgency(ctx context.Context, agency *domain.Agency) error {
	query := `UPDATE agency SET name=$1,address=$2,reg_id=$3,updated_at=$4 WHERE agency_id=$5;`
	row := h.db.QueryRowContext(ctx, query, agency.Name, agency.Address, agency.RegistrationID,agency.UpdatedAt,agency.AgencyID)
	err := row.Err()
	if err != nil {
		return err
	}
	return nil
}

func (h *agencyRepositoryDB) DeleteAgency(ctx context.Context, agencyID uuid.UUID) error {
	query := `DELETE FROM agency WHERE agency_id=$1;`
	_, err := h.db.ExecContext(ctx, query, agencyID)
	if err != nil {
		return err
	}
	return nil
}
