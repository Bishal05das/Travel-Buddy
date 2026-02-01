package mocks

import (

	"github.com/bishal05das/travelbuddy/internal/domain"
)

type MockTourRepository struct {
	tours map[int]*domain.Tour
	err error
}
func NewMockTourRepository() *MockTourRepository{
	return &MockTourRepository{
		tours: map[int]*domain.Tour{},
		
	}
}

func (m *MockTourRepository) CreateTour(tour *domain.Tour) error {
	if m.err != nil {
		return m.err
	}
	tour.ID=len(m.tours)+1
	m.tours[tour.ID]=tour
	return nil
}

func (m *MockTourRepository) ListTour(agencyID int) ([]*domain.Tour)