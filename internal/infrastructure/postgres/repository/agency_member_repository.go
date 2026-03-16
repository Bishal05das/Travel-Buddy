package repository

import (
	"context"
	"database/sql"

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
	query := `INSERT INTO agency_members (agency_id,role_id,name,email,phone,password) VALUES ($1,$2,$3,$4,$5,$6) RETURNING member_id;`

	return h.db.QueryRowContext(ctx, query, member.AgencyID,member.RoleID,member.Name,member.Email,member.Phone,member.Password).Scan(&member.MemberID)
}

func (h *agencyMemberRepositoryDB) ListMember(ctx context.Context, agencyID uuid.UUID) ([]*domain.ListMemberResponse, error) {
	// var members []*domain.ListMemberResponse
	// query := `SELECT name,email,phone,password FROM agency_members WHERE agency_id=$1;`
	// err := h.db.SelectContext(ctx, &members, query, agencyID)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return nil, nil
	// 	}
	// 	return nil, err
	// }
	// return members, nil

	query := `
	SELECT 
		m.member_id,
		m.name,
		m.email,
		m.phone,
		p.permission_id
	FROM agency_members m
	LEFT JOIN role_permissions rp ON m.role_id = rp.role_id
	LEFT JOIN permissions p ON rp.permission_id = p.permission_id
	WHERE m.agency_id = $1
	ORDER BY m.member_id;
`
	rows,err := h.db.QueryxContext(ctx, query, agencyID)
	if err != nil {
		return nil,err
	}
	defer rows.Close()

	memberMap := map[uuid.UUID]*domain.ListMemberResponse{}
	for rows.Next() {
		var (
			memberID uuid.UUID
			perm domain.Permission
		)

		var m domain.ListMemberResponse
		err := rows.Scan(&memberID,&m.Name,&m.Email,&m.Phone,&perm.PermissionID)
		if err != nil {
			return nil, err
		}
		if _, exists := memberMap[memberID]; !exists {
			m.MemberID = memberID
			m.Permissions = []int{}
			memberMap[memberID]=&m
		}

		if perm.PermissionID != 0 {
			memberMap[memberID].Permissions = append(memberMap[memberID].Permissions, perm.PermissionID)
		}
	}

	var members []*domain.ListMemberResponse
	for _,v := range memberMap {
		members = append(members, v)
	}
	return members, nil

}

func (h *agencyMemberRepositoryDB) UpdateMember(ctx context.Context,member *domain.AgencyMember) error {
	query := `UPDATE agency_members SET name=$1,email=$2,phone=$3,password=$4 WHERE member_id=$5;`
	_,err := h.db.ExecContext(ctx, query, member.Name,member.Email,member.Phone,member.Password,member.MemberID)
	if err != nil {
		return err
	}
	return nil
}

func (h *agencyMemberRepositoryDB) DeleteMember(ctx context.Context, memberID uuid.UUID) error {
	query := `DELETE FROM agency_members WHERE member_id=$1;`
	_, err := h.db.ExecContext(ctx, query, memberID)
	if err != nil {
		return err
	}
	return nil
}


func (h *agencyMemberRepositoryDB) GetRoleIDFromMemberIDForUpdate(ctx context.Context,memberID uuid.UUID) (*int,error) {
	var roleID int

	query := `SELECT role_id FROM agency_members WHERE member_id = $1 FOR UPDATE;`
	err := h.executor(ctx).QueryRowxContext(ctx,query,memberID).Scan(&roleID)
	if err != nil {
		return nil,err
	}
	return &roleID,nil
}

func (h *agencyMemberRepositoryDB) FindMember(ctx context.Context,email string) (*domain.AgencyMember,error) {
	var member domain.AgencyMember
	query := `SELECT member_id,agency_id,name,password,phone,role_id FROM agency_members WHERE email=$1;`
	err := h.db.GetContext(ctx,&member,query,email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil,nil
		}
		return nil,err
	}
	return &member,nil
}

func (h *agencyMemberRepositoryDB) executor(ctx context.Context) sqlx.ExtContext {
	if tx, ok := GetTx(ctx); ok {
		return tx
	}
	return h.db
}