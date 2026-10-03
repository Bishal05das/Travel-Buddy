package domain

import "github.com/google/uuid"

type AgencyMember struct {
	MemberID uuid.UUID `json:"member_id" db:"member_id"`
	AgencyID uuid.UUID `json:"agency_id" db:"agency_id"`
	RoleID   int       `json:"role_id" db:"role_id"`
	Name     string    `json:"name" db:"name"`
	Email    string    `json:"email" db:"email"`
	Phone    string    `json:"phone" db:"phone"`
	Password string    `json:"-" db:"password"`
	IsOwner  bool      `json:"is_owner" db:"is_owner"`
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
	RoleName    string    `json:"role_name" db:"role_name"`
	IsOwner     bool      `json:"is_owner" db:"is_owner"`
	Permissions []int     `json:"permissions"`
}

type MemberProfile struct {
	MemberID   uuid.UUID `json:"member_id" db:"member_id"`
	AgencyID   uuid.UUID `json:"agency_id" db:"agency_id"`
	AgencyName string    `json:"agency_name" db:"agency_name"`
	Name       string    `json:"name" db:"name"`
	Email      string    `json:"email" db:"email"`
	Phone      string    `json:"phone" db:"phone"`
	RoleName   string    `json:"role_name" db:"role_name"`
	IsOwner    bool      `json:"is_owner" db:"is_owner"`
}

type UpdatePermissionRequest struct {
	Permissions []int `json:"permissions" validate:"required,dive,gt=0"`
}
