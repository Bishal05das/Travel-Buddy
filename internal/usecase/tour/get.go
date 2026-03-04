package tourusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type GetTourUseCase struct {
	repo port.TourRepository
}

func NewGetTourUseCase(repo port.TourRepository) *GetTourUseCase {
	return &GetTourUseCase{
		repo: repo,
	}
}

func(uc *GetTourUseCase) Execute(ctx context.Context,tourID uuid.UUID) (*domain.Tour,error) {
	return uc.repo.GetByID(ctx,tourID)
}