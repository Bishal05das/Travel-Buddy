package userusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/config"
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
	"golang.org/x/crypto/bcrypt"
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
	usr, err := uc.userRepo.FindUserByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, domain.ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(user.Password))
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}
	accessToken, err := util.CreateJWT(uc.cnf.JWTSecretkey, util.Payload{
		UserID: usr.UserID,
		Role:   usr.Role,
	}, uc.cnf.JWTTTL)
	if err != nil {
		return nil, err
	}
	return &accessToken, nil
}
