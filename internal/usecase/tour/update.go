package tourusecase

import (
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type UpdateTourUseCase struct {
	repo port.TourRepository
}

func NewUpdateTourUseCase(repo port.TourRepository) *UpdateTourUseCase {
	return &UpdateTourUseCase{
		repo: repo,
	}
}

func (uc *UpdateTourUseCase)Execute(tour *domain.Tour) error {

	return uc.repo.UpdateTour(tour)
}