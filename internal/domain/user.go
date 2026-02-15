package domain

import "github.com/google/uuid"

type User struct {
	UserID   uuid.UUID
	Name     string
	Email    string
	Password string
	Phone    int
	Role     string
}

type ReqLogin struct {
	Email    string
	Password string
}
