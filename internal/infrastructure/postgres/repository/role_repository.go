package repository

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type roleRepositoryDB struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) port.RoleRepository {
	return &roleRepositoryDB{
		db: db,
	}
}

func (r *roleRepositoryDB) CreateRole(ctx context.Context, role *domain.Role) error {
	query := `INSERT INTO roles (agency_id,role_name) VALUES ($1,$2) RETURNING role_id;`

	return r.db.QueryRowContext(ctx, query, role.AgencyID, role.RoleName).Scan(&role.RoleID)
}

func (r *roleRepositoryDB) DeleteRole(ctx context.Context, roleID int) error {
	query := `DELETE FROM roles WHERE role_id=$1;`
	_, err := r.db.ExecContext(ctx, query, roleID)
	if err != nil {
		return err
	}
	return nil
}

// func (r *roleRepositoryDB) CreateRolePermission(ctx context.Context, roleID int, permissionID int) error {
// 	query := `INSERT INTO role_permissions (role_id,permission_id) VALUES ($1,$2) ON CONFLICT (role_id, permission_id) DO NOTHING;`

// 	_, err := r.db.ExecContext(ctx, query, roleID, permissionID)
// 	return err
// }

func (r *roleRepositoryDB) DeletePermissionsFromRole(ctx context.Context, roleID int) error {
	query := `DELETE FROM role_permissions WHERE role_id=$1;`
	_, err := r.db.ExecContext(ctx, query, roleID)
	return err
}

func (r *roleRepositoryDB) AddPermissionsToRole(ctx context.Context, roleID int, permissionIDs []int) error {

	const q = `INSERT INTO role_permissions (role_id, permission_id) 
	SELECT $1, unnest($2::int[])
	ON CONFLICT (role_id, permission_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, q, roleID, pq.Array(permissionIDs))
	return err
}
