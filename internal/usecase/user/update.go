package userusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type UpdateUserUseCase struct {
	repo port.UserRepository
}

func NewUpdateUserUseCase(repo port.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		repo: repo,
	}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context, user *domain.User) error {

	return uc.repo.UpdateUser(ctx, user)
}
