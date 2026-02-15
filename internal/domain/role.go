package domain

import "github.com/google/uuid"

type Role struct {
	RoleID       int
	AgencyID     uuid.UUID
	RoleName     string
	Description  string
	IsSystemRole bool
}
