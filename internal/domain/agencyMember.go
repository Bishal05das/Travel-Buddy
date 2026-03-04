package domain

import "github.com/google/uuid"

type AgencyMember struct {
	MemberID uuid.UUID `json:"member_id" db:"member_id"`
	AgencyID uuid.UUID `json:"agency_id" db:"agency_id"`
	RoleID   int       `json:"role_id" db:"role_id"`
	Name     string    `json:"name" db:"name"`
	Email    string    `json:"email" db:"email"`
	Phone    string    `json:"phone" db:"phone"`
	Password string    `json:"password" db:"password"`
}

type CreateMemberRequest struct {
	AgencyID    uuid.UUID `json:"agency_id" db:"agency_id"`
	Name        string    `json:"name" db:"name"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	Password    string    `json:"password" db:"password"`
	RoleName    string    `json:"role_name"`
	Permissions []int     `json:"permissions" db:"permissions"`
}

type ListMemberResponse struct {
	MemberID    uuid.UUID `json:"member_id" db:"member_id"`
	Name        string    `json:"name" db:"name"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	Permissions []int     `json:"permissions"`
}

type UpdatePermissionRequest struct {
	Permissions []int
}
