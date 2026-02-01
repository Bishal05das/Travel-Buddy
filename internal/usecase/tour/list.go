package tourusecase

import "github.com/bishal05das/travelbuddy/internal/domain"

type ListTourUseCase struct {
	repo domain.TourRepository
}

func NewListTourUseCase(repo domain.TourRepository) *ListTourUseCase {
	return &ListTourUseCase{
		repo: repo,
	}
}

func(uc *ListTourUseCase) Execute(agencyID int) ([]*domain.Tour,error) {
	return uc.repo.ListTour(agencyID)
}