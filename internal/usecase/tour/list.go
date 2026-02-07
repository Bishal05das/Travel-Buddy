package tourusecase

import (
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type ListTourUseCase struct {
	repo port.TourRepository
}

func NewListTourUseCase(repo port.TourRepository) *ListTourUseCase {
	return &ListTourUseCase{
		repo: repo,
	}
}

func(uc *ListTourUseCase) Execute(agencyID int) ([]*domain.Tour,error) {
	return uc.repo.ListTour(agencyID)
}