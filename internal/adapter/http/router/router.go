package router

import (
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
	middleware "github.com/bishal05das/travelbuddy/internal/adapter/http/middlewares"
)

type Router struct {
	mux               *http.ServeMux
	middleware        *middleware.Middleware
	homeHandler       *handler.HomeHandler
	searchHandler     *handler.SearchHandler
	tourHandler       *handler.TourHandler
	userHandler       *handler.UserHandler
	bookingHandler    *handler.BookingHandler
	agencyHandler     *handler.AgencyHandler
	memberHandler     *handler.MemberHandler
	permissionHandler *handler.PermissionHandler
}

func NewRoutes(
	mux *http.ServeMux,
	middleware *middleware.Middleware,
	homeHandler *handler.HomeHandler,
	searchHandler     *handler.SearchHandler,
	tourHandler *handler.TourHandler,
	userHandler *handler.UserHandler,
	bookingHandler *handler.BookingHandler,
	agencyHandler *handler.AgencyHandler,
	memberHandler *handler.MemberHandler,
	permissionHandler *handler.PermissionHandler) *Router {

	return &Router{
		mux:               mux,
		middleware:        middleware,
		homeHandler:       homeHandler,
		searchHandler: searchHandler,
		tourHandler:       tourHandler,
		userHandler:       userHandler,
		bookingHandler:    bookingHandler,
		agencyHandler:     agencyHandler,
		memberHandler:     memberHandler,
		permissionHandler: permissionHandler,
	}
}

func (r *Router) RegisterRoutes() {
	//home
	r.mux.HandleFunc("GET /home", r.homeHandler.GetHome)

	//search
	r.mux.HandleFunc("GET /search",r.searchHandler.Search)
	//tour
	r.mux.HandleFunc("POST /tours", r.tourHandler.Create)
	r.mux.HandleFunc("GET /tours/{tour_id}", r.tourHandler.Get)
	r.mux.HandleFunc("GET /tours/list/{agency_id}", r.tourHandler.List)
	r.mux.HandleFunc("PUT /tours/{tour_id}", r.tourHandler.Update)
	r.mux.HandleFunc("DELETE /tours/{tour_id}", r.tourHandler.Delete)

	//user
	r.mux.HandleFunc("POST /users", r.userHandler.CreateUser)
	r.mux.HandleFunc("POST /users/login", r.userHandler.UserLogin)
	r.mux.HandleFunc("DELETE /users/{user_id}", r.userHandler.DeleteUser)
	r.mux.HandleFunc("PUT /users/{user_id}", r.userHandler.UpdateUser)

	//booking
	r.mux.HandleFunc("POST /bookings/{tour_id}", r.bookingHandler.CreateBooking)

	//agency
	r.mux.HandleFunc("POST /agency", r.agencyHandler.CreateAgency)
	r.mux.HandleFunc("PUT /agency", r.agencyHandler.UpdateAgency)
	r.mux.HandleFunc("DELETE /agency/{id}", r.agencyHandler.DeleteAgency)

	//member
	r.mux.HandleFunc("POST /members", r.memberHandler.CreateMember)
	r.mux.HandleFunc("DELETE /members/{member_id}", r.memberHandler.DeleteMember)
	r.mux.HandleFunc("GET /members/{agency_id}", r.memberHandler.ListMember)
	r.mux.HandleFunc("PUT /members/{member_id}/permissions", r.memberHandler.UpdateMemberPermissions)
	r.mux.HandleFunc("POST /members/login", r.memberHandler.MemberLogin)

	//permissions
	r.mux.HandleFunc("POST /permissions", r.permissionHandler.CreatePermission)
	r.mux.HandleFunc("DELETE /permissions/{id}", r.permissionHandler.DeletePermission)

}
