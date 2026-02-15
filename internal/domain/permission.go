package domain

type Permission struct {
	PermissionID int
	Name         string
	Resource     string
	Action       string
}
