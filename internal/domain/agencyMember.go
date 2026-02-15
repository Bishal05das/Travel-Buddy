package domain

import "github.com/google/uuid"

type AgencyMember struct {
	MemberID uuid.UUID
	AgencyID uuid.UUID
	RoleID   int
	Name     string
	Email    string
	Phone    int
	Password string
	Role     Role
}
