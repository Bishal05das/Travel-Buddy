package userusecase

import (
	"context"
	"errors"

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
	usr, err := uc.repo.FindUserByID(ctx, user.UserID)
	if err != nil {
		return err
	}
	if usr == nil {
		return errors.New("User Not Found")
	}

	return uc.repo.UpdateUser(ctx, user)
}
