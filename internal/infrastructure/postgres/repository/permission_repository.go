package repository

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type permissionRepositoryDB struct {
	db *sqlx.DB
}

func NewPermissionRepositoryDB(db *sqlx.DB) *permissionRepositoryDB {
	return &permissionRepositoryDB{
		db: db,
	}
}

func (h *permissionRepositoryDB) CreatePermission(ctx context.Context, permisson *domain.Permission) error {
	query := `INSERT INTO permissions (name,resource,action) VALUES ($1,$2,$3) RETURNING permission_id;`

	return h.executor(ctx).QueryRowxContext(ctx, query, permisson.Name, permisson.Resource, permisson.Action).Scan(&permisson.PermissionID)
}

func (h *permissionRepositoryDB) DeletePermission(ctx context.Context, permissionID int) error {
	query := `DELETE FROM permissions WHERE permission_id=$1;`
	res, err := h.executor(ctx).ExecContext(ctx, query, permissionID)
	if err != nil {
		return err
	}
	return requireAffected(res, "permission not found")
}

func (h *permissionRepositoryDB) ListPermissions(ctx context.Context) ([]domain.Permission, error) {
	perms := []domain.Permission{}
	query := `SELECT permission_id, name, resource, action FROM permissions ORDER BY resource, action`
	if err := sqlx.SelectContext(ctx, h.executor(ctx), &perms, query); err != nil {
		return nil, err
	}
	return perms, nil
}

// MemberHasPermission reports whether the member belongs to agencyID and
// they own the agency or their role grants resource:action. It reads current data on every call,
// so revoked permissions and deleted members take effect immediately.
func (h *permissionRepositoryDB) MemberHasPermission(ctx context.Context, memberID, agencyID uuid.UUID, resource, action string) (bool, error) {
	query := `
	SELECT EXISTS (
		SELECT 1
		FROM agency_members m
		WHERE m.member_id = $1 AND m.agency_id = $2
		  AND (m.is_owner OR EXISTS (
			SELECT 1 FROM role_permissions rp
			JOIN permissions p ON p.permission_id=rp.permission_id
			WHERE rp.role_id=m.role_id AND p.resource=$3 AND p.action=$4
		  ))
	);`
	var ok bool
	err := h.executor(ctx).QueryRowxContext(ctx, query, memberID, agencyID, resource, action).Scan(&ok)
	return ok, err
}

func (h *permissionRepositoryDB) executor(ctx context.Context) sqlx.ExtContext {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return h.db
}
