// internal/domain/agency_image.go
package domain

import (
	"github.com/google/uuid"
	"time"
)

type AgencyImage struct {
	ImageID    uuid.UUID  `db:"image_id"`
	AgencyID   uuid.UUID  `db:"agency_id"`
	ImagePath  string     `db:"image_path"`
	IsActive   bool       `db:"is_active"`
	UploadedBy *uuid.UUID `db:"uploaded_by"`
	CreatedAt  time.Time  `db:"created_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}
