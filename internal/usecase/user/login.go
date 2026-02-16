package userusecase

import (
	"context"
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

func NewUserLoginUseCase(userRepo port.UserRepository, cnf *config.Config) port.LoginUser {
	return &userLoginUseCase{
		userRepo: userRepo,
		cnf:      cnf,
	}
}

func (uc *userLoginUseCase) Execute(ctx context.Context, user *domain.ReqLogin) (*string, error) {
	usr, err := uc.userRepo.FindUser(ctx, user.Email, user.Password)
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
