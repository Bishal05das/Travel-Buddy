package tourusecase

import (
	"context"
	"errors"

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
	hasTour, err := uc.repo.GetByID(ctx,tour.TourID)
	if err != nil {
		return err
	}
	if hasTour == nil {
		return errors.New("Tour Not Found")
	}

	return uc.repo.UpdateTour(ctx, tour)
}
