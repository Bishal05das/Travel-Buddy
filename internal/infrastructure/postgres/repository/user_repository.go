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

type userRepositoryDB struct {
	db *sqlx.DB
}

func NewUserRepositoryDB(db *sqlx.DB) port.UserRepository {
	return &userRepositoryDB{
		db: db,
	}
}

func (h *userRepositoryDB) CreateUser(ctx context.Context,user *domain.User) error {

	query := `INSERT INTO users (name,email,password,phone) VALUES ($1,$2,$3,$4) RETURNING user_id;`

	return h.db.QueryRowContext(ctx,query, user.Name, user.Email, user.Password, user.Phone).Scan(&user.UserID)
}

func (h *userRepositoryDB) UpdateUser(ctx context.Context,user *domain.User) error {
	query := `UPDATE users SET name=$1,email=$2,password=$3,phone=$4,updated_at=$5 WHERE user_id=$6;`
	res,err := h.db.ExecContext(ctx,query, user.Name, user.Email, user.Password, user.Phone, user.UpdatedAt, user.UserID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (h *userRepositoryDB) DeleteUser(ctx context.Context,UserID uuid.UUID) error {
	query := `DELETE FROM users WHERE user_id=$1;`
	_, err := h.db.ExecContext(ctx,query, UserID)
	if err != nil {
		return err
	}
	return nil
}

func (h *userRepositoryDB) FindUserByEmail(ctx context.Context,email string) (*domain.User,error) {
	var user domain.User
	query := `SELECT user_id,name,phone,role FROM users WHERE email=$1`
	err := h.db.GetContext(ctx,&user,query,email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil,nil
		}
		return nil,err
	}
	return &user,nil
}

func (h *userRepositoryDB) FindUserByID(ctx context.Context,id uuid.UUID) (*domain.User,error) {
	var user domain.User
	query := `SELECT name,email,phone,role FROM users WHERE user_id=$1`
	err := h.db.GetContext(ctx,&user,query,id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil,nil
		}
		return nil,err
	}
	return &user,nil
}
