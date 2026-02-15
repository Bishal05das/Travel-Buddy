package repository

import (
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
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

func (h *userRepositoryDB) CreateUser(user *domain.User) error {
	query := `INSERT INTO users (name,email,password,phone) VALUES ($1,$2,$3,$4) RETURNING id;`

	return h.db.QueryRow(query, user.Name, user.Email, user.Password, user.Phone).Scan(&user.ID)
}

func (h *userRepositoryDB) UpdateUser(user *domain.User) error {
	query := `UPDATE users SET name=$1,email=$2,password=$3,phone=$4,role=$5;`
	row := h.db.QueryRow(query, user.Name, user.Email, user.Password, user.Phone, user.Role)
	err := row.Err()
	if err != nil {
		return err
	}
	return nil
}

func (h *userRepositoryDB) DeleteUser(UserID int) error {
	query := `DELETE FROM users WHERE id=$1;`
	_, err := h.db.Exec(query, UserID)
	if err != nil {
		return err
	}
	return nil
}
