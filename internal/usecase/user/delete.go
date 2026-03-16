package userusecase

import (
	"context"
	"errors"

	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type DeleteUserUseCase struct {
	repo port.UserRepository
}

func NewDeleteUserUseCase(repo port.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		repo: repo,
	}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, userID uuid.UUID) error {
	usr, err := uc.repo.FindUserByID(ctx,userID)
	if err != nil {
		return err
	}
	if usr == nil {
		return errors.New("User Not Found")
	}
	return uc.repo.DeleteUser(ctx,userID)
}
