package port

import "github.com/bishal05das/travelbuddy/internal/domain"

type CreateTour interface {
	Execute(tour *domain.Tour) error
}
