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
}

type CreateMemberRequest struct {
	AgencyID    uuid.UUID
	Name        string
	Email       string
	Phone       int
	Password    string
	RoleName    string
	Permissions []int
}

type UpdatePermissionRequest struct {
	MemberID    uuid.UUID
	Permissions []int
}
