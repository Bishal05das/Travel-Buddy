package tourusecase

import (
	"fmt"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
)

type CreateTourUseCase struct {
	repo domain.TourRepository
}

func NewCreateTourUseCase(r domain.TourRepository) *CreateTourUseCase {
	return &CreateTourUseCase{
		repo: r,
	}
}

func (uc *CreateTourUseCase) Execute(tour *domain.Tour) error {
	//add business logic here
	if tour.AgencyID <= 0 ||
		tour.Name == "" ||
		tour.Description == "" ||
		tour.EndDate.IsZero() ||
		tour.StartDate.IsZero() ||
		tour.LastEnrollmentDate.IsZero() ||
		tour.Price <= 0 {
		return fmt.Errorf("invalid or missing tour data")
	}
	if tour.StartDate.Before(time.Now()) || tour.EndDate.Before(time.Now()) {
		return fmt.Errorf("enter correct start date or end date")
	}
	if tour.EndDate.Before(tour.StartDate) {
		return fmt.Errorf("End Date can not be before start Date")
	}
	return uc.repo.CreateTour(tour)
}
