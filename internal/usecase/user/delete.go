package userusecase

import (
	"context"

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
	return uc.repo.DeleteUser(ctx,userID)
}
