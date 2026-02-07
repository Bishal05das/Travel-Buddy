package userusecase



import (
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

func (uc *CreateUserUseCase) Execute(user *domain.User) error {
	//add business logic here

	return uc.repo.CreateUser(user)
}
