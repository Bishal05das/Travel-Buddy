package userusecase

import (
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

func(uc *DeleteUserUseCase) Execute(userID uuid.UUID) error {
return uc.repo.DeleteUser(userID)
}