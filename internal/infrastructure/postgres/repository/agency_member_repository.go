package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type agencyMemberRepositoryDB struct {
	db *sqlx.DB
}

func NewAgencyMemberRepositoryDB(db *sqlx.DB) port.AgencyMemberRepository {
	return &agencyMemberRepositoryDB{
		db: db,
	}
}

func (h *agencyMemberRepositoryDB) CreateMember(ctx context.Context, member *domain.AgencyMember) error {
	query := `INSERT INTO agency_members (agency_id,role_id,name,email,phone,password,is_owner) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING member_id;`

	return h.executor(ctx).QueryRowxContext(ctx, query, member.AgencyID, member.RoleID, member.Name, member.Email, member.Phone, member.Password, member.IsOwner).Scan(&member.MemberID)
}

func (h *agencyMemberRepositoryDB) ListMember(ctx context.Context, agencyID uuid.UUID) ([]*domain.ListMemberResponse, error) {
	query := `
	SELECT 
		m.member_id,
		m.name,
		m.email,
		m.phone,
		r.role_name,
		m.is_owner,
		p.permission_id
	FROM agency_members m
	JOIN roles r ON r.role_id = m.role_id
	LEFT JOIN permissions p ON m.is_owner OR EXISTS (
		SELECT 1 FROM role_permissions rp
		WHERE rp.role_id = m.role_id AND rp.permission_id = p.permission_id
	)
	WHERE m.agency_id = $1
	ORDER BY m.name, m.member_id, p.permission_id;
`
	rows, err := h.db.QueryxContext(ctx, query, agencyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []*domain.ListMemberResponse{}
	memberMap := map[uuid.UUID]*domain.ListMemberResponse{}
	for rows.Next() {
		var (
			m            domain.ListMemberResponse
			permissionID sql.NullInt64 // NULL for members whose role has no permissions
		)
		if err := rows.Scan(&m.MemberID, &m.Name, &m.Email, &m.Phone, &m.RoleName, &m.IsOwner, &permissionID); err != nil {
			return nil, err
		}
		member, exists := memberMap[m.MemberID]
		if !exists {
			m.Permissions = []int{}
			member = &m
			memberMap[m.MemberID] = member
			members = append(members, member)
		}
		if permissionID.Valid {
			member.Permissions = append(member.Permissions, int(permissionID.Int64))
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return members, nil
}

func (h *agencyMemberRepositoryDB) UpdateMember(ctx context.Context, member *domain.AgencyMember) error {
	query := `UPDATE agency_members SET name=$1,email=$2,phone=$3,password=$4 WHERE member_id=$5;`
	_, err := h.db.ExecContext(ctx, query, member.Name, member.Email, member.Phone, member.Password, member.MemberID)
	if err != nil {
		return err
	}
	return nil
}

func (h *agencyMemberRepositoryDB) DeleteMember(ctx context.Context, memberID uuid.UUID, agencyScope *uuid.UUID) error {
	// The ownership predicate applies to every caller, including platform admins.
	query := `DELETE FROM agency_members
		WHERE member_id=$1 AND ($2::uuid IS NULL OR agency_id=$2)
		AND NOT is_owner RETURNING member_id;`
	var deletedID uuid.UUID
	err := h.db.QueryRowxContext(ctx, query, memberID, agencyScope).Scan(&deletedID)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	var isOwner bool
	err = h.db.GetContext(ctx, &isOwner,
		`SELECT is_owner FROM agency_members WHERE member_id=$1 AND ($2::uuid IS NULL OR agency_id=$2)`,
		memberID, agencyScope)
	if err == nil && isOwner {
		return domain.ErrOwnerProtected
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return errors.New("member not found")
}

func (h *agencyMemberRepositoryDB) GetRoleIDFromMemberIDForUpdate(ctx context.Context, memberID uuid.UUID, agencyScope *uuid.UUID) (*int, error) {
	var roleID int
	var ownerRole bool

	// A shared role must not offer a way to change an owner's permissions.
	query := `SELECT m.role_id, EXISTS (
		SELECT 1 FROM agency_members owner WHERE owner.role_id=m.role_id AND owner.is_owner
	) FROM agency_members m
	WHERE m.member_id=$1 AND ($2::uuid IS NULL OR m.agency_id=$2) FOR UPDATE OF m;`
	err := h.executor(ctx).QueryRowxContext(ctx, query, memberID, agencyScope).Scan(&roleID, &ownerRole)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("member not found")
	}
	if err != nil {
		return nil, err
	}
	if ownerRole {
		return nil, domain.ErrOwnerProtected
	}
	return &roleID, nil
}

func (h *agencyMemberRepositoryDB) GetPermissionIDs(ctx context.Context, memberID uuid.UUID) ([]int, error) {
	query := `
	SELECT p.permission_id
	FROM agency_members m
	JOIN permissions p ON m.is_owner OR EXISTS (
		SELECT 1 FROM role_permissions rp
		WHERE rp.role_id=m.role_id AND rp.permission_id=p.permission_id
	)
	WHERE m.member_id = $1;`
	ids := []int{}
	if err := sqlx.SelectContext(ctx, h.executor(ctx), &ids, query, memberID); err != nil {
		return nil, err
	}
	return ids, nil
}

func (h *agencyMemberRepositoryDB) GetMemberProfile(ctx context.Context, memberID, agencyID uuid.UUID) (*domain.MemberProfile, error) {
	var profile domain.MemberProfile
	query := `SELECT m.member_id, m.agency_id, a.name AS agency_name,
		m.name, m.email, m.phone, r.role_name, m.is_owner
		FROM agency_members m
		JOIN agency a ON a.agency_id=m.agency_id
		JOIN roles r ON r.role_id=m.role_id
		WHERE m.member_id=$1 AND m.agency_id=$2`
	if err := h.db.GetContext(ctx, &profile, query, memberID, agencyID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrMemberNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func (h *agencyMemberRepositoryDB) FindMember(ctx context.Context, email string) (*domain.AgencyMember, error) {
	var member domain.AgencyMember
	query := `SELECT member_id,agency_id,name,password,phone,role_id FROM agency_members WHERE email=$1;`
	err := h.db.GetContext(ctx, &member, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &member, nil
}

func (h *agencyMemberRepositoryDB) executor(ctx context.Context) sqlx.ExtContext {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return h.db
}
