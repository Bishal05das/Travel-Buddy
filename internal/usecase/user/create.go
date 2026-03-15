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
	//add business logic here
	if user.Name == "" {
		return errors.New("name is required")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}

	if user.Password == "" {
		return errors.New("password is required")
	}

	if user.Phone == "" {
		return errors.New("phone is required")
	}
	hashedPassword,err := util.HashPassword(user.Password)
	if err != nil {
		return errors.New("error in password hashing")
	}
	user.Password = hashedPassword
	return uc.repo.CreateUser(ctx, user)
}
