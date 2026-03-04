package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

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

	//UseCase
	homeUC := homeusecase.NewHomeUseCase(homeRepo)

	createTourUC := tourusecase.NewCreateTourUseCase(tourRepo)
	getTourUC := tourusecase.NewGetTourUseCase(tourRepo)
	listTourUC := tourusecase.NewListTourUseCase(tourRepo)
	deleteTourUC := tourusecase.NewDeleteTourUseCase(tourRepo)
	updateTourUC := tourusecase.NewUpdateTourUseCase(tourRepo)

	createuserUC := userusecase.NewCreateUserUseCase(userRepo)
	createBookingUC := bookingusecase.NewCreateBookingUseCase(txManager, bookingRepo, tourRepo, paymentRepo)
	loginUserUC := userusecase.NewUserLoginUseCase(userRepo, cfg)
	deleteUserUC := userusecase.NewDeleteUserUseCase(userRepo)
	updateUserUC := userusecase.NewUpdateUserUseCase(userRepo)

	createAgencyUC := agencyusecase.NewCreateAgencyUseCase(agencyRepo)
	deleteAgencyUC := agencyusecase.NewDeleteAgencyUseCase(agencyRepo)
	updateAgencyUC := agencyusecase.NewUpdateAgencyUseCase(agencyRepo)

	createMemberUC := memberusecase.NewCreateAgencyMemberUseCase(txManager, memberRepo, roleRepo)
	deleteMemberUC := memberusecase.NewDeleteAgencyMemberUseCase(memberRepo)
	listMemberUC := memberusecase.NewListAgencyMemberUseCase(memberRepo)
	LoginMemberUC := memberusecase.NewMemberLoginUseCase(memberRepo,cfg)

	updatePermissionUC := memberusecase.NewUpdatePermissionUseCase(txManager, memberRepo, roleRepo)
	createPermissionsUC := permissionusecase.NewCreatePermissionUseCase(permissionRepo)
	deletePermissionUC := permissionusecase.NewDeletePermissionUseCase(permissionRepo)

	//handler
	homeHandler := handler.NewHomeHandler(homeUC)
	tourHandler := handler.NewTourHandler(createTourUC,getTourUC,listTourUC,updateTourUC,deleteTourUC)
	userHandler := handler.NewUserHandler(createuserUC, loginUserUC,deleteUserUC,updateUserUC)
	bookingHandler := handler.NewBookingHandler(createBookingUC)
	agencyHandler := handler.NewAgencyHandler(createAgencyUC, updateAgencyUC, deleteAgencyUC)
	memberHandler := handler.NewMemberHandler(createMemberUC, deleteMemberUC, listMemberUC, updatePermissionUC,LoginMemberUC)
	permissionHandler := handler.NewPermissionHandler(createPermissionsUC, deletePermissionUC)

	//middleware
	middleware := middleware.NewMiddleware(cfg)

	//router setup
	router := router.NewRoutes(mux, middleware,homeHandler, tourHandler, userHandler, bookingHandler, agencyHandler, memberHandler, permissionHandler)
	router.RegisterRoutes()

	fmt.Println("Listening to server on port 3000")
	addr := ":" + strconv.Itoa(cfg.HttpPort)
	err = http.ListenAndServe(addr, newHandler)
	if err != nil {
		fmt.Println(fmt.Println("Server failed to start:", err))
	}
}
