package router

import (
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
)

type Router struct {
	mux *http.ServeMux
	tourHandler *handler.TourHandler
	userHandler *handler.UserHandler
	bookingHandler *handler.BookingHandler
}
func NewRoutes(mux *http.ServeMux,tourHandler *handler.TourHandler,userHandler *handler.UserHandler,bookingHandler *handler.BookingHandler) *Router {
	return &Router{
		mux: mux,
		tourHandler: tourHandler,
		userHandler: userHandler,
		bookingHandler: bookingHandler,
	}
}

func(r *Router) RegisterRoutes(){
	r.mux.HandleFunc("POST /tours",r.tourHandler.Create)
	r.mux.HandleFunc("POST /users",r.userHandler.CreateUser)
	r.mux.HandleFunc("POST /bookings",r.bookingHandler.CreateBooking)
}
