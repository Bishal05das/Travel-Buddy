package router

import (
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
	middleware "github.com/bishal05das/travelbuddy/internal/adapter/http/middlewares"
)

type Router struct {
	mux               *http.ServeMux
	middleware        *middleware.MiddlewareManager
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
	middleware *middleware.MiddlewareManager,
	homeHandler *handler.HomeHandler,
	searchHandler *handler.SearchHandler,
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
		searchHandler:     searchHandler,
		tourHandler:       tourHandler,
		userHandler:       userHandler,
		bookingHandler:    bookingHandler,
		agencyHandler:     agencyHandler,
		memberHandler:     memberHandler,
		permissionHandler: permissionHandler,
	}
}

func (r *Router) public(h http.Handler) http.Handler {
	return middleware.Chain(
		h,
		r.middleware.Logger,
		r.middleware.RateLimiter,
	)
}

func (r *Router) protected(h http.Handler) http.Handler {
	return middleware.Chain(
		h,
		r.middleware.Logger,
		r.middleware.RateLimiter,
		r.middleware.Authentication,
	)
}

func (r *Router) RegisterRoutes() {

	// HOME
	r.mux.Handle(
		"GET /home",
		r.public(http.HandlerFunc(r.homeHandler.GetHome)),
	)

	// SEARCH
	r.mux.Handle(
		"GET /search",
		r.public(http.HandlerFunc(r.searchHandler.Search)),
	)

	// TOURS
	r.mux.Handle(
		"POST /tours",
		r.protected(http.HandlerFunc(r.tourHandler.Create)),
	)

	r.mux.Handle(
		"GET /tours/{tour_id}",
		r.public(http.HandlerFunc(r.tourHandler.Get)),
	)

	r.mux.Handle(
		"GET /tours/list/{agency_id}",
		r.public(http.HandlerFunc(r.tourHandler.List)),
	)

	r.mux.Handle(
		"PUT /tours/{tour_id}",
		r.protected(http.HandlerFunc(r.tourHandler.Update)),
	)

	r.mux.Handle(
		"DELETE /tours/{tour_id}",
		r.protected(http.HandlerFunc(r.tourHandler.Delete)),
	)

	// USERS
	r.mux.Handle(
		"POST /users",
		r.public(http.HandlerFunc(r.userHandler.CreateUser)),
	)

	r.mux.Handle(
		"POST /users/login",
		r.public(http.HandlerFunc(r.userHandler.UserLogin)),
	)

	r.mux.Handle(
		"DELETE /users/{user_id}",
		r.protected(http.HandlerFunc(r.userHandler.DeleteUser)),
	)

	r.mux.Handle(
		"PUT /users/{user_id}",
		r.protected(http.HandlerFunc(r.userHandler.UpdateUser)),
	)

	// BOOKINGS
	r.mux.Handle(
		"POST /bookings/{tour_id}",
		r.protected(http.HandlerFunc(r.bookingHandler.CreateBooking)),
	)

	// AGENCY
	r.mux.Handle(
		"POST /agency",
		r.protected(http.HandlerFunc(r.agencyHandler.CreateAgency)),
	)

	r.mux.Handle(
		"PUT /agency/{agency_id}",
		r.protected(http.HandlerFunc(r.agencyHandler.UpdateAgency)),
	)

	r.mux.Handle(
		"DELETE /agency/{agency_id}",
		r.protected(http.HandlerFunc(r.agencyHandler.DeleteAgency)),
	)

	// MEMBERS
	r.mux.Handle(
		"POST /members",
		r.protected(http.HandlerFunc(r.memberHandler.CreateMember)),
	)

	r.mux.Handle(
		"DELETE /members/{member_id}",
		r.protected(http.HandlerFunc(r.memberHandler.DeleteMember)),
	)

	r.mux.Handle(
		"GET /members/{agency_id}",
		r.protected(http.HandlerFunc(r.memberHandler.ListMember)),
	)

	r.mux.Handle(
		"PUT /members/{member_id}/permissions",
		r.protected(http.HandlerFunc(r.memberHandler.UpdateMemberPermissions)),
	)

	r.mux.Handle(
		"POST /members/login",
		r.public(http.HandlerFunc(r.memberHandler.MemberLogin)),
	)

	// PERMISSIONS
	r.mux.Handle(
		"POST /permissions",
		r.protected(http.HandlerFunc(r.permissionHandler.CreatePermission)),
	)

	r.mux.Handle(
		"DELETE /permissions/{id}",
		r.protected(http.HandlerFunc(r.permissionHandler.DeletePermission)),
	)
}
