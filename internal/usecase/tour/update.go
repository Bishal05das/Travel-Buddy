package tourusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type updateTourUseCase struct {
	repo port.TourRepository
}

func NewUpdateTourUseCase(repo port.TourRepository) port.UpdateTour {
	return &updateTourUseCase{
		repo: repo,
	}
}

func (uc *updateTourUseCase) Execute(ctx context.Context, tour *domain.Tour) error {

	return uc.repo.UpdateTour(ctx, tour)
}
