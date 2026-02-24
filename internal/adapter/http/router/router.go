package router

import (
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
	middleware "github.com/bishal05das/travelbuddy/internal/adapter/http/middlewares"
)

type Router struct {
	mux               *http.ServeMux
	middleware        *middleware.Middleware
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
	tourHandler *handler.TourHandler,
	userHandler *handler.UserHandler,
	bookingHandler *handler.BookingHandler,
	agencyHandler *handler.AgencyHandler,
	memberHandler *handler.MemberHandler,
	permissionHandler *handler.PermissionHandler) *Router {

	return &Router{
		mux:               mux,
		middleware:        middleware,
		tourHandler:       tourHandler,
		userHandler:       userHandler,
		bookingHandler:    bookingHandler,
		agencyHandler:     agencyHandler,
		memberHandler:     memberHandler,
		permissionHandler: permissionHandler,
	}
}

func (r *Router) RegisterRoutes() {
	//tour
	r.mux.HandleFunc("POST /tours", r.tourHandler.Create)

	//user
	r.mux.Handle("POST /users", http.HandlerFunc(r.userHandler.CreateUser))
	r.mux.Handle("POST /users/login", http.HandlerFunc(r.userHandler.UserLogin))

	//booking
	r.mux.HandleFunc("POST /bookings", r.bookingHandler.CreateBooking)

	//agency
	r.mux.HandleFunc("POST /agency", r.agencyHandler.CreateAgency)
	r.mux.HandleFunc("PUT /agency", r.agencyHandler.UpdateAgency)
	r.mux.HandleFunc("DELETE /agency/{id}", r.agencyHandler.DeleteAgency)

	//member
	r.mux.HandleFunc("POST /member", r.memberHandler.CreateMember)
	r.mux.HandleFunc("PUT /member/permissions", r.memberHandler.UpdateMemberPermissions)
	r.mux.HandleFunc("DELETE /member", r.memberHandler.DeleteMember)

	//permissions
	r.mux.HandleFunc("POST /permissions", r.permissionHandler.CreatePermission)
	r.mux.HandleFunc("DELETE /permissions/{id}", r.permissionHandler.DeletePermission)
}
