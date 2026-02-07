package userusecase

import (
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type DeleteUserUseCase struct {
	repo port.UserRepository
}

func NewDeleteUserUseCase(repo port.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		repo: repo,
	}
}

func(uc *DeleteUserUseCase) Execute(userID int) error {
return uc.repo.DeleteUser(userID)
}