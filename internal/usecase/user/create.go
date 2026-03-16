package userusecase

import (
	"context"
	"errors"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
)

type CreateUserUseCase struct {
	repo port.UserRepository
}

func NewCreateUserUseCase(r port.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		repo: r,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, user *domain.User) error {
	usr, err := uc.repo.FindUserByEmail(ctx,user.Email)
	if err != nil {
		return err
	}
	if usr != nil {
		return errors.New("email already exist")
	}

	hashedPassword,err := util.HashPassword(user.Password)
	if err != nil {
		return errors.New("error in password hashing")
	}
	user.Password = hashedPassword
	return uc.repo.CreateUser(ctx, user)
}

