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
	AgencyID    uuid.UUID `json:"agency_id" validate:"required,uuid"`
	Name        string    `json:"name" validate:"required,min=3,max=120"`
	Email       string    `json:"email" validate:"required,email"`
	Phone       string    `json:"phone" validate:"required,e164"`
	Password    string    `json:"password" validate:"required,min=8,max=64"`
	RoleName    string    `json:"role_name" validate:"required,min=2,max=100"`
	Permissions []int     `json:"permissions" validate:"required,min=1,dive,gt=0"`
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
