package tourusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type createTourUseCase struct {
	repo port.TourRepository
}

func NewCreateTourUseCase(r port.TourRepository) port.CreateTour {
	return &createTourUseCase{
		repo: r,
	}
}

func (uc *createTourUseCase) Execute(ctx context.Context, tour *domain.Tour) error {


	return uc.repo.CreateTour(ctx, tour)
}
