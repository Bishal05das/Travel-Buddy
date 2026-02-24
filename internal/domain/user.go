package domain

import "github.com/google/uuid"

type User struct {
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Name     string    `json:"name" db:"name"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"password" db:"password"`
	Phone    string    `json:"phone" db:"phone"`
	Role     string    `json:"role" db:"role"`
}

type ReqLogin struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
