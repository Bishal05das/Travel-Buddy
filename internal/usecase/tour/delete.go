package tourusecase

import (
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type DeleteTourUseCase struct {
	repo port.TourRepository
}

func NewDeleteTourUseCase(repo port.TourRepository) *DeleteTourUseCase {
	return &DeleteTourUseCase{
		repo: repo,
	}
}

func(uc *DeleteTourUseCase) Execute(tourID int) error {
return uc.repo.DeleteTour(tourID)
}