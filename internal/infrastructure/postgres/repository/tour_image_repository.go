package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/google/uuid"
)

// Lock the tour before replacing its image so concurrent replacements cannot
// create two active images. The agency filter also prevents cross-agency edits.
func (h *tourRepositoryDB) UpdateTourImage(ctx context.Context, agencyID, tourID uuid.UUID, imagePath string) (string, error) {
	tx, err := h.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	var status string
	err = tx.QueryRowContext(ctx, `SELECT status FROM tours WHERE tour_id=$1 AND agency_id=$2 FOR UPDATE`, tourID, agencyID).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		return "", domain.ErrImageTargetNotFound
	}
	if err != nil {
		return "", err
	}
	if status == "cancelled" {
		return "", domain.ErrTourCancelled
	}

	var oldPath string
	err = tx.QueryRowContext(ctx, `UPDATE tour_images SET is_active=FALSE, deleted_at=CURRENT_TIMESTAMP WHERE tour_id=$1 AND is_active=TRUE AND deleted_at IS NULL RETURNING image_path`, tourID).Scan(&oldPath)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO tour_images (tour_id, image_path, is_active) VALUES ($1, $2, TRUE)`, tourID, imagePath); err != nil {
		return "", err
	}
	if err := tx.Commit(); err != nil {
		return "", err
	}
	return oldPath, nil
}
