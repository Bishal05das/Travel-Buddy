package tourusecase

import (
	"context"
	"errors"
	"github.com/bishal05das/travelbuddy/internal/domain"
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

func (uc *updateTourStatusUseCase) Execute(ctx context.Context, actor domain.Actor, tourID uuid.UUID, status string) error {
	switch status {
	case "open", "closed", "cancelled":
	default:
		return errors.New("status must be one of: open, closed, cancelled")
	}
	return uc.repo.UpdateTourStatus(ctx, tourID, status, actor.AgencyScope())
}
