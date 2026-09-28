package tourusecase

import (
	"context"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type updateTourStatusUseCase struct {
	repo port.TourRepository
}

func NewUpdateTourStatusUseCase(repo port.TourRepository) port.UpdateTourStatus {
	return &updateTourStatusUseCase{
		repo: repo,
	}
}

func (uc *updateTourStatusUseCase) Execute(ctx context.Context, tourID uuid.UUID, status string) error {
	return uc.repo.UpdateTourStatus(ctx, tourID, status)
}
