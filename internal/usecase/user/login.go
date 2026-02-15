package userusecase

import (
	"errors"

	"github.com/bishal05das/travelbuddy/config"
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
)

type userLoginUseCase struct {
	userRepo port.UserRepository
	cnf      *config.Config
}

func NewUserLoginUseCase(userRepo port.UserRepository, cnf *config.Config) 

func (uc *userLoginUseCase) Execute(user *domain.ReqLogin) (*string, error) {
	usr, err := uc.userRepo.FindUser(user.Email, user.Password)
	if usr == nil {
		// http.Error(w, "Invalid Credentials", http.StatusBadRequest)
		return nil, errors.New("Invalid Credentials")
	}
	accessToken, err := util.CreateJWT(uc.cnf.JWTSecretkey, util.Payload{
		UserID: usr.UserID,
		Role:   usr.Role,
	})
	if err != nil {
		//http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, err
	}
	return &accessToken, nil
}
