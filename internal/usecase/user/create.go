package userusecase

import (
	"context"


	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
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

	return uc.repo.CreateUser(ctx, user)
}
