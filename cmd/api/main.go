package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bishal05das/travelbuddy/config"
	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
	middleware "github.com/bishal05das/travelbuddy/internal/adapter/http/middlewares"
	"github.com/bishal05das/travelbuddy/internal/adapter/http/router"
	db "github.com/bishal05das/travelbuddy/internal/infrastructure/postgres"
	"github.com/bishal05das/travelbuddy/internal/infrastructure/postgres/repository"
	bookingusecase "github.com/bishal05das/travelbuddy/internal/usecase/booking"
	tourusecase "github.com/bishal05das/travelbuddy/internal/usecase/tour"
	userusecase "github.com/bishal05das/travelbuddy/internal/usecase/user"
)

func main() {
	mux := http.NewServeMux()
	cfg := config.GetConfig()

	dbCon, err := db.NewConnection(cfg)
	if err != nil {
		fmt.Println("err in database connection: ", err)
		return
	}
	err = db.MigrateDB(dbCon, cfg)
	if err != nil {
		fmt.Println("DB Migration failed:", err)
		os.Exit(1)
	}
	newHandler := middleware.Cors(mux)
	tourRepo := repository.NewTourRepositoryDB(dbCon)
	createtouruc := tourusecase.NewCreateTourUseCase(tourRepo)
	tourHandler := handler.NewTourHandler(createtouruc)
	userRepo := repository.NewUserRepositoryDB(dbCon)
	createuseruc := userusecase.NewCreateUserUseCase(userRepo)
	loginUC := userusecase.NewUserLoginUseCase(userRepo,cfg)
	userHandler := handler.NewUserHandler(createuseruc,loginUC)
	bookingRepo := repository.NewBookingRepository(dbCon)
	paymentRepo := repository.NewPaymentRepositoryDB(dbCon)
	createBookingUC := bookingusecase.NewCreateBookingUseCase(bookingRepo,tourRepo,paymentRepo)
	bookingHandler := handler.NewBookingHandler(createBookingUC)
	router := router.NewRoutes(mux, tourHandler,userHandler,bookingHandler)
	router.RegisterRoutes()
	fmt.Println("Listening to server on port 3000")
	err = http.ListenAndServe(":3000", newHandler)
	if err != nil {
		fmt.Println(fmt.Println("Server failed to start:", err))
	}
}
