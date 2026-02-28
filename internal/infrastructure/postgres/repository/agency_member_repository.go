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
	var members []*domain.ListMemberResponse
	query := `SELECT name,email,phone,password FROM agency_members WHERE agency_id=$1;`
	err := h.db.SelectContext(ctx, &members, query, agencyID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
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


func (h *agencyMemberRepositoryDB) GetRoleIDFromMemberID(ctx context.Context,memberID uuid.UUID) (*int,error) {
	var roleID int

	query := `SELECT role_id FROM agency_members WHERE member_id = $1`
	err := h.db.QueryRowContext(ctx,query,memberID).Scan(&roleID)
	if err != nil {
		return nil,err
	}
	return &roleID,nil
}