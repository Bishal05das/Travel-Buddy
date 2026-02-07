package router

import (
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
)

type Routes struct {
	mux         *http.ServeMux
	tourHandler *handler.TourHandler
}

func NewRoutes(mux *http.ServeMux, tourHandler *handler.TourHandler) *Routes {
	return &Routes{
		mux:         mux,
		tourHandler: tourHandler,
	}
}

func (r *Routes) RegisterRoutes() {
	r.mux.HandleFunc("POST /tours", r.tourHandler.Create)
}
