package domain

type Role string

const (
	RoleAdmin      Role = "ADMIN"
	RoleUser       Role = "USER"
	RoleSuperAdmin Role = "SUPER"
	RoleSubAdmin   Role = "SUB"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	Phone    int
	Role     Role
}
