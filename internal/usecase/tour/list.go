package tourusecase

import (
	"context"


	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type listTourUseCase struct {
	repo port.TourRepository
}

func NewListTourUseCase(repo port.TourRepository) port.ListTour {
	return &listTourUseCase{
		repo: repo,
	}
}

func (uc *listTourUseCase) Execute(ctx context.Context, agencyID uuid.UUID) ([]*domain.Tour, error) {
	return uc.repo.ListTour(ctx, agencyID)
}
