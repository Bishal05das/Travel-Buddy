package main

import (
	"fmt"
	"net/http"

	"github.com/bishal05das/travelbuddy/config"
	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
	"github.com/bishal05das/travelbuddy/internal/adapter/http/router"
	"github.com/bishal05das/travelbuddy/internal/repository"
	tourusecase "github.com/bishal05das/travelbuddy/internal/usecase/tour"
	"github.com/bishal05das/travelbuddy/pkg/db"
)

func main(){
	mux := http.NewServeMux()
	cfg := config.GetConfig()

	db,err := db.NewConnection(cfg)
	if err != nil {
		fmt.Println("err in database connection: ",err)
		return
	}
	tourRepo := repository.NewTourRepositoryDB(db)
	tourusecase := tourusecase.NewCreateTourUseCase(tourRepo)
	tourHandler := handler.NewTourHandler(tourusecase)
	router := router.NewRoutes(mux,tourHandler)
	router.RegisterRoutes()
	fmt.Println("Listening to server on port 3000")
	err = http.ListenAndServe(":3000",mux)
	if err != nil {
		fmt.Println(fmt.Println("Server failed to start:", err))
	}
}