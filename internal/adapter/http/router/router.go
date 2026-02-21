package router

import (
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
)

type Router struct {
	mux            *http.ServeMux
	tourHandler    *handler.TourHandler
	userHandler    *handler.UserHandler
	bookingHandler *handler.BookingHandler
	agencyHandler  *handler.AgencyHandler
	memberHandler  *handler.MemberHandler
}

func NewRoutes(mux *http.ServeMux, tourHandler *handler.TourHandler, userHandler *handler.UserHandler, bookingHandler *handler.BookingHandler, agencyHandler *handler.AgencyHandler, memberHandler *handler.MemberHandler) *Router {
	return &Router{
		mux:            mux,
		tourHandler:    tourHandler,
		userHandler:    userHandler,
		bookingHandler: bookingHandler,
		agencyHandler:  agencyHandler,
		memberHandler:  memberHandler,
	}
}

func (r *Router) RegisterRoutes() {
	r.mux.HandleFunc("POST /tours", r.tourHandler.Create)
	r.mux.HandleFunc("POST /users", r.userHandler.CreateUser)
	r.mux.HandleFunc("POST /bookings", r.bookingHandler.CreateBooking)
	r.mux.HandleFunc("POST /agency", r.agencyHandler.CreateAgency)
	r.mux.HandleFunc("PUT /agency",r.agencyHandler.UpdateAgency)
	r.mux.HandleFunc("DELETE /agency",r.agencyHandler.DeleteAgency)
	r.mux.HandleFunc("POST /member",r.memberHandler.CreateMember)
	r.mux.HandleFunc("PUT /member/permissions",r.memberHandler.UpdateMemberPermissions)
	r.mux.HandleFunc("DELETE /member",r.memberHandler.DeleteMember)
}
