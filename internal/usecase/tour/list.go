package tourusecase

import (
	"context"


	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type ListTourUseCase struct {
	repo port.TourRepository
}

func NewListTourUseCase(repo port.TourRepository) *ListTourUseCase {
	return &ListTourUseCase{
		repo: repo,
	}
}

func (uc *ListTourUseCase) Execute(ctx context.Context, agencyID uuid.UUID) ([]*domain.Tour, error) {
	return uc.repo.ListTour(ctx, agencyID)
}
