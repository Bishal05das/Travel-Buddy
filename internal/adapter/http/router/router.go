package router

import (
	"net/http"
	"strings"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
	middleware "github.com/bishal05das/travelbuddy/internal/adapter/http/middlewares"
	"github.com/bishal05das/travelbuddy/internal/domain"
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

// protected requires a valid token, then applies the given authorization
// checks in order.
func (r *Router) protected(h http.Handler, authz ...middleware.Middleware) http.Handler {
	chain := append([]middleware.Middleware{
		r.middleware.Logger,
		r.middleware.RateLimiter,
		r.middleware.Authentication,
	}, authz...)
	return middleware.Chain(h, chain...)
}

func (r *Router) RegisterRoutes() {
	m := r.middleware
	superOnly := m.RequireRole(domain.RoleSuper)

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
		"POST /agency/{agency_id}/tours",
		r.protected(http.HandlerFunc(r.tourHandler.Create), m.RequirePermission("tour", "create")),
	)

	r.mux.Handle(
		"GET /tours/{tour_id}",
		r.public(http.HandlerFunc(r.tourHandler.Get)),
	)

	r.mux.Handle(
		"GET /agency/{agency_id}/tours/list",
		r.public(http.HandlerFunc(r.tourHandler.List)),
	)

	r.mux.Handle(
		"PUT /agency/{agency_id}/tours/{tour_id}",
		r.protected(http.HandlerFunc(r.tourHandler.Update), m.RequirePermission("tour", "update")),
	)

	r.mux.Handle(
		"PATCH /tours/{tour_id}/tour-status",
		r.protected(http.HandlerFunc(r.tourHandler.UpdateStatus), m.RequirePermission("tour", "update")),
	)

	r.mux.Handle(
		"DELETE /tours/{tour_id}",
		r.protected(http.HandlerFunc(r.tourHandler.Delete), m.RequirePermission("tour", "delete")),
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
		r.protected(http.HandlerFunc(r.userHandler.DeleteUser), m.RequireSelfOrSuper("user_id")),
	)

	r.mux.Handle(
		"PUT /users/{user_id}",
		r.protected(http.HandlerFunc(r.userHandler.UpdateUser), m.RequireSelfOrSuper("user_id")),
	)

	// BOOKINGS
	r.mux.Handle(
		"POST /bookings/{tour_id}",
		r.protected(http.HandlerFunc(r.bookingHandler.CreateBookingByUser), m.RequireRole(domain.RoleUser)),
	)
	r.mux.Handle(
		"POST /admin/bookings/{tour_id}",
		r.protected(http.HandlerFunc(r.bookingHandler.CreateBookingByAdmin), m.RequireRole(domain.RoleMember), m.RequirePermission("booking", "create")),
	)

	// AGENCY
	r.mux.Handle(
		"POST /agency",
		r.protected(http.HandlerFunc(r.agencyHandler.CreateAgency), superOnly),
	)

	r.mux.Handle(
		"PUT /agency/{agency_id}",
		r.protected(http.HandlerFunc(r.agencyHandler.UpdateAgency), m.RequirePermission("agency", "update")),
	)

	r.mux.Handle(
		"DELETE /agency/{agency_id}",
		r.protected(http.HandlerFunc(r.agencyHandler.DeleteAgency), m.RequirePermission("agency", "delete")),
	)

	// MEMBERS
	r.mux.Handle(
		"POST /members/{agency_id}",
		r.protected(http.HandlerFunc(r.memberHandler.CreateMember), m.RequirePermission("member", "create")),
	)

	r.mux.Handle(
		"DELETE /members/{member_id}",
		r.protected(http.HandlerFunc(r.memberHandler.DeleteMember), m.RequirePermission("member", "delete")),
	)

	r.mux.Handle(
		"GET /members/{agency_id}",
		r.protected(http.HandlerFunc(r.memberHandler.ListMember), m.RequirePermission("member", "read")),
	)

	r.mux.Handle(
		"PUT /members/{member_id}/permissions",
		r.protected(http.HandlerFunc(r.memberHandler.UpdateMemberPermissions), m.RequirePermission("member", "update")),
	)

	r.mux.Handle(
		"POST /members/login",
		r.public(http.HandlerFunc(r.memberHandler.MemberLogin)),
	)

	// PERMISSIONS
	r.mux.Handle(
		"POST /permissions",
		r.protected(http.HandlerFunc(r.permissionHandler.CreatePermission), superOnly),
	)

	r.mux.Handle(
		"DELETE /permissions/{id}",
		r.protected(http.HandlerFunc(r.permissionHandler.DeletePermission), superOnly),
	)

	// images
	imageFS := http.FileServer(http.Dir("./images"))
	r.mux.Handle(
		"GET /images/", http.StripPrefix("/images/", serveFilesOnly(imageFS)))
}

// serveFilesOnly blocks directory listings, which would expose every
// uploaded file name, and stops browsers from sniffing content types.
func serveFilesOnly(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" || strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		h.ServeHTTP(w, r)
	})
}
