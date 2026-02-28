package domain

type Permission struct {
	PermissionID int    `json:"permission_id" db:"permission_id"`
	Name         string `json:"name" db:"name"`
	Resource     string `json:"resource" db:"resource"`
	Action       string `json:"action" db:"action"`
}
