package userusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type UpdateUserUseCase struct {
	repo port.UserRepository
}

func NewUpdateUserUseCase(repo port.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		repo: repo,
	}
}

func (uc *UpdateUserUseCase) Execute(ctx context.Context,userID uuid.UUID, user *domain.UpdateUserReq) error {
	updatedUser := &domain.User{
		UserID: userID,
		Name: user.Name,
		Email: user.Email,
		Password: user.Password,
		Phone: user.Phone,
	}

	return uc.repo.UpdateUser(ctx, updatedUser)
}
