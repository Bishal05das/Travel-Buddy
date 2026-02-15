package router

import (
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
)

type Router struct {
	mux *http.ServeMux
	tourHandler *handler.TourHandler
}
func NewRoutes(mux *http.ServeMux,tourHandler *handler.TourHandler) *Router {
	return &Router{
		mux: mux,
		tourHandler: tourHandler,
	}
}

func(r *Router) RegisterRoutes(){
	r.mux.HandleFunc("POST /tours",r.tourHandler.Create)
}
