package repository

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/jmoiron/sqlx"
)

type permissionRepositoryDB struct {
	db *sqlx.DB
}

func NewPermissionRepositoryDB(db *sqlx.DB) *permissionRepositoryDB{
	return &permissionRepositoryDB{
		db: db,
	}
}

func (h *permissionRepositoryDB) CreatePermission(ctx context.Context, permisson *domain.Permission) error {
	query := `INSERT INTO permissions (name,resource,action) VALUES ($1,$2,$3) RETURNING id;`

	return h.db.QueryRowContext(ctx, query, permisson.Name, permisson.Resource, permisson.Action).Scan(&permisson.PermissionID)
}

func (h *permissionRepositoryDB) DeletePermission(ctx context.Context,permissionID int) error {
	query := `DELETE FROM permissions WHERE id=$1;`
	_,err := h.db.ExecContext(ctx,query,permissionID)
	if err != nil {
		return err
	}
	return nil
}
