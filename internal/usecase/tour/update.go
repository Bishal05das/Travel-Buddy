package tourusecase

import "github.com/bishal05das/travelbuddy/internal/domain"

type UpdateTourUseCase struct {
	repo domain.TourRepository
}

func NewUpdateTourUseCase(repo domain.TourRepository) *UpdateTourUseCase {
	return &UpdateTourUseCase{
		repo: repo,
	}
}

func (uc *UpdateTourUseCase)Execute(tour *domain.Tour) error {
	
	return uc.repo.UpdateTour(tour)
}