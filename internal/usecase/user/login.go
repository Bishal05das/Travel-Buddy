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
	usr, err := uc.userRepo.FindUser(ctx, user.Email)
	if usr == nil {
		return nil, errors.New("Invalid Credentials")
	}
	reqHashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, errors.New("error in hashing password")
	}
	if reqHashedPassword != usr.Password {
		return nil, errors.New("Invalid Password")
	}  
	accessToken, err := util.CreateJWT(uc.cnf.JWTSecretkey, util.Payload{
		UserID: usr.UserID,
		Role:   usr.Role,
	})
	if err != nil {
		return nil, err
	}
	return &accessToken, nil
}
