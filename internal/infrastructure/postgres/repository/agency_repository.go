package repository

import (
	"context"
	"database/sql"
	"fmt"

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

func (h *agencyRepositoryDB) CreateAgency(ctx context.Context,agency *domain.Agency,imagePath string) error {
	tx, err := h.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// 1) Insert agency
	const agencyQ = `
        INSERT INTO agency (name, address, reg_id)
        VALUES ($1, $2, $3)
        RETURNING agency_id;
    `
	if err := tx.QueryRowContext(
		ctx, agencyQ,
		agency.Name, agency.Address, agency.RegistrationID,
	).Scan(&agency.AgencyID); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("insert agency: %w", err)
	}

	// 2) Insert first image as active
	const imageQ = `
        INSERT INTO agency_images (agency_id, image_path, is_active)
        VALUES ($1, $2, TRUE);
    `
	if _, err := tx.ExecContext(ctx, imageQ, agency.AgencyID, imagePath); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("insert agency image: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// UpdateAgencyImage marks current image inactive and inserts a new active one.
// Returns the OLD image path so the handler can delete the file from disk.
func (h *agencyRepositoryDB) UpdateAgencyImage(
	ctx context.Context,
	agencyID uuid.UUID,
	newImagePath string,
) (string, error) {
	tx, err := h.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("begin tx: %w", err)
	}

	var oldPath string
	const deactivateQ = `
        UPDATE agency_images
        SET is_active = FALSE,
            deleted_at = CURRENT_TIMESTAMP
        WHERE agency_id = $1 AND is_active = TRUE
        RETURNING image_path;
    `
	err = tx.QueryRowContext(ctx, deactivateQ, agencyID).Scan(&oldPath)
	if err != nil && err != sql.ErrNoRows {
		_ = tx.Rollback()
		return "", fmt.Errorf("deactivate old image: %w", err)
	}

	const insertQ = `
        INSERT INTO agency_images (agency_id, image_path, is_active)
        VALUES ($1, $2, TRUE);
    `
	if _, err := tx.ExecContext(ctx, insertQ, agencyID, newImagePath); err != nil {
		_ = tx.Rollback()
		return "", fmt.Errorf("insert new image: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit: %w", err)
	}
	return oldPath, nil
}

func (h *agencyRepositoryDB) GetCurrentImage(
	ctx context.Context,
	agencyID uuid.UUID,
) (*domain.AgencyImage, error) {
	const q = `
        SELECT image_id, agency_id, image_path, is_active,
               uploaded_by, created_at, deleted_at
        FROM agency_images
        WHERE agency_id = $1 AND is_active = TRUE AND deleted_at IS NULL
        LIMIT 1;
    `
	var img domain.AgencyImage
	if err := h.db.GetContext(ctx, &img, q, agencyID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &img, nil
}

func (h *agencyRepositoryDB) ListAgencyImages(
	ctx context.Context,
	agencyID uuid.UUID,
) ([]*domain.AgencyImage, error) {
	const q = `
        SELECT image_id, agency_id, image_path, is_active,
               uploaded_by, created_at, deleted_at
        FROM agency_images
        WHERE agency_id = $1
        ORDER BY created_at DESC;
    `
	var imgs []*domain.AgencyImage
	if err := h.db.SelectContext(ctx, &imgs, q, agencyID); err != nil {
		return nil, err
	}
	return imgs, nil
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
	row := h.db.QueryRowContext(ctx, query, agency.Name, agency.Address, agency.RegistrationID, agency.UpdatedAt, agency.AgencyID)
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
