package tourusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type deleteTourUseCase struct {
	repo port.TourRepository
}

func NewDeleteTourUseCase(repo port.TourRepository) port.DeleteTour {
	return &deleteTourUseCase{
		repo: repo,
	}
}

func (uc *deleteTourUseCase) Execute(ctx context.Context, tourID uuid.UUID) error {
	return uc.repo.DeleteTour(ctx, tourID)
}
