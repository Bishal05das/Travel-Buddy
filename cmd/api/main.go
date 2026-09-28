package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/bishal05das/travelbuddy/config"
	"github.com/bishal05das/travelbuddy/internal/adapter/http/handler"
	middleware "github.com/bishal05das/travelbuddy/internal/adapter/http/middlewares"
	"github.com/bishal05das/travelbuddy/internal/adapter/http/router"
	db "github.com/bishal05das/travelbuddy/internal/infrastructure/postgres"
	"github.com/bishal05das/travelbuddy/internal/infrastructure/postgres/repository"
	agencyusecase "github.com/bishal05das/travelbuddy/internal/usecase/agency"
	memberusecase "github.com/bishal05das/travelbuddy/internal/usecase/agencyMember"
	bookingusecase "github.com/bishal05das/travelbuddy/internal/usecase/booking"
	homeusecase "github.com/bishal05das/travelbuddy/internal/usecase/home"
	permissionusecase "github.com/bishal05das/travelbuddy/internal/usecase/permission"
	searchusecase "github.com/bishal05das/travelbuddy/internal/usecase/search"
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
	defer dbCon.Close()

	// 10 MB multipart uploads plus room for the other form fields.
	newHandler := middleware.Cors(middleware.LimitBody(11 << 20)(mux))
	txManager := repository.NewTxManager(dbCon)

	//Repository

	tourRepo := repository.NewTourRepositoryDB(dbCon)
	userRepo := repository.NewUserRepositoryDB(dbCon)
	bookingRepo := repository.NewBookingRepository(dbCon)
	paymentRepo := repository.NewPaymentRepositoryDB(dbCon)
	agencyRepo := repository.NewAgencyRepositoryDB(dbCon)
	memberRepo := repository.NewAgencyMemberRepositoryDB(dbCon)
	roleRepo := repository.NewRoleRepository(dbCon)
	permissionRepo := repository.NewPermissionRepositoryDB(dbCon)
	homeRepo := repository.NewHomeRepositoryDB(dbCon)
	searchRepo := repository.NewSearchRepository(dbCon)

	//UseCase
	homeUC := homeusecase.NewHomeUseCase(homeRepo)

	searchUC := searchusecase.NewSearchUseCase(searchRepo)

	createTourUC := tourusecase.NewCreateTourUseCase(tourRepo)
	getTourUC := tourusecase.NewGetTourUseCase(tourRepo)
	listTourUC := tourusecase.NewListTourUseCase(tourRepo)
	deleteTourUC := tourusecase.NewDeleteTourUseCase(tourRepo)
	updateTourUC := tourusecase.NewUpdateTourUseCase(tourRepo)
	updateTourStatusUC := tourusecase.NewUpdateTourStatusUseCase(tourRepo)

	createuserUC := userusecase.NewCreateUserUseCase(userRepo)
	createBookingUC := bookingusecase.NewCreateBookingUseCase(txManager, bookingRepo, tourRepo, paymentRepo)
	listBookingsUC := bookingusecase.NewListBookingsUseCase(bookingRepo)
	getBookingUC := bookingusecase.NewGetBookingUseCase(bookingRepo)
	updateBookingStatusUC := bookingusecase.NewUpdateBookingStatusUseCase(txManager, bookingRepo, tourRepo, paymentRepo)
	cancelMyBookingUC := bookingusecase.NewCancelMyBookingUseCase(txManager, bookingRepo, tourRepo, paymentRepo)
	loginUserUC := userusecase.NewUserLoginUseCase(userRepo, cfg)
	deleteUserUC := userusecase.NewDeleteUserUseCase(userRepo)
	updateUserUC := userusecase.NewUpdateUserUseCase(userRepo)

	createAgencyUC := agencyusecase.NewCreateAgencyUseCase(agencyRepo)
	deleteAgencyUC := agencyusecase.NewDeleteAgencyUseCase(agencyRepo)
	updateAgencyUC := agencyusecase.NewUpdateAgencyUseCase(agencyRepo)

	createMemberUC := memberusecase.NewCreateAgencyMemberUseCase(txManager, memberRepo, roleRepo)
	deleteMemberUC := memberusecase.NewDeleteAgencyMemberUseCase(memberRepo)
	listMemberUC := memberusecase.NewListAgencyMemberUseCase(memberRepo)
	LoginMemberUC := memberusecase.NewMemberLoginUseCase(memberRepo, cfg)

	updatePermissionUC := memberusecase.NewUpdatePermissionUseCase(txManager, memberRepo, roleRepo)
	createPermissionsUC := permissionusecase.NewCreatePermissionUseCase(permissionRepo)
	deletePermissionUC := permissionusecase.NewDeletePermissionUseCase(permissionRepo)

	//handler
	homeHandler := handler.NewHomeHandler(homeUC)
	searchHandler := handler.NewSearchHandler(searchUC)
	tourHandler := handler.NewTourHandler(createTourUC, getTourUC, listTourUC, updateTourUC, updateTourStatusUC, deleteTourUC)
	userHandler := handler.NewUserHandler(createuserUC, loginUserUC, deleteUserUC, updateUserUC)
	bookingHandler := handler.NewBookingHandler(createBookingUC, listBookingsUC, getBookingUC, updateBookingStatusUC, cancelMyBookingUC)
	agencyHandler := handler.NewAgencyHandler(createAgencyUC, updateAgencyUC, deleteAgencyUC)
	memberHandler := handler.NewMemberHandler(createMemberUC, deleteMemberUC, listMemberUC, updatePermissionUC, LoginMemberUC)
	permissionHandler := handler.NewPermissionHandler(createPermissionsUC, deletePermissionUC)

	//middleware
	middleware := middleware.NewMiddlewareManager(cfg, permissionRepo)

	//router setup
	router := router.NewRoutes(mux, middleware, homeHandler, searchHandler, tourHandler, userHandler, bookingHandler, agencyHandler, memberHandler, permissionHandler)
	router.RegisterRoutes()

	srv := &http.Server{
		Addr:              ":" + strconv.Itoa(cfg.HttpPort),
		Handler:           newHandler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)
	go func() {
		fmt.Println("Listening to server on port ", cfg.HttpPort)
		serverErr <- srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server failed to start:", err)
			os.Exit(1)
		}
	case <-stop:
		fmt.Println("Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Println("Graceful shutdown failed:", err)
		}
	}
}
