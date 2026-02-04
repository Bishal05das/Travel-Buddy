package tourusecase

import "github.com/bishal05das/travelbuddy/internal/domain"

type DeleteTourUseCase struct {
	repo domain.TourRepository
}

func NewDeleteTourUseCase(repo domain.TourRepository) *DeleteTourUseCase {
	return &DeleteTourUseCase{
		repo: repo,
	}
}

func(uc *DeleteTourUseCase) Execute(tourID int) error {
return uc.repo.DeleteTour(tourID)
}