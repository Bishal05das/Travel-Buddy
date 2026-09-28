package userusecase

import (
	"context"
	"errors"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
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

	other, err := uc.repo.FindUserByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if other != nil && other.UserID != user.UserID {
		return domain.ErrEmailTaken
	}

	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return errors.New("error in password hashing")
	}
	user.Password = hashedPassword

	return uc.repo.UpdateUser(ctx, user)
}
