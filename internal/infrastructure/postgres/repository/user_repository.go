package repository

import (
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryDB struct {
	db *sqlx.DB
}

func NewUserRepositoryDB(db *sqlx.DB) *UserRepositoryDB {
	return &UserRepositoryDB{
		db: db,
	}
}

func (h *UserRepositoryDB) CreateUser(user *domain.User) error {
	query := `INSERT INTO users (name,email,password,phone) VALUES ($1,$2,$3,$4) RETURNING id;`

	return h.db.QueryRow(query, user.Name, user.Email, user.Password, user.Phone).Scan(&user.ID)
}

func (h *UserRepositoryDB) UpdateUser(user *domain.User) error {
	query := `UPDATE users SET name=$1,email=$2,password=$3,phone=$4,role=$5;`
	row := h.db.QueryRow(query, user.Name, user.Email, user.Password, user.Phone, user.Role)
	err := row.Err()
	if err != nil {
		return err
	}
	return nil
}

func (h *UserRepositoryDB) DeleteUser(UserID int) error {
	query := `DELETE FROM users WHERE id=$1;`
	_, err := h.db.Exec(query, UserID)
	if err != nil {
		return err
	}
	return nil
}
