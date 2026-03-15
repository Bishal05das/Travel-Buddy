package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	RoleID     int
	AgencyID   uuid.UUID
	RoleName   string
	Created_at time.Time
}
