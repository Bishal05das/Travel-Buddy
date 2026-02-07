package domain

type Permission string

const (
	PermCreateTour   Permission = "CREATE_TOUR"
	PermUpdateTour   Permission = "UPDATE_TOUR"
	PermDeleteTour   Permission = "DELETE_TOUR"
	PermViewBooking  Permission = "VIEW_BOOKING"
	PermCancelBooking Permission = "CANCEL_BOOKING"
)

type AgencyMember struct {
	ID int
	AgencyID int
	Name string
	Email string
	Password string
	Role Role
	Permissions []Permission
}