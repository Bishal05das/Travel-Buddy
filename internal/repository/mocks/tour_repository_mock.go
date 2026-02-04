package mocks

import (
	"errors"

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

func (m *MockTourRepository) ListTour(agencyID int) ([]*domain.Tour,error) {
	var result []*domain.Tour

	for _,tour := range m.tours {
		if tour.AgencyID == agencyID {
			result = append(result, tour)
		}
	}
	return result,nil
}

func (m *MockTourRepository) UpdateTour(t *domain.Tour) error {

	existing, ok := m.tours[t.ID]
	if !ok {
		return errors.New("tour not found")
	}
	existing.AgencyID=t.AgencyID
	existing.Name=t.Name
	existing.StartDate=t.StartDate
	existing.EndDate=t.EndDate
	existing.Description=t.Description
	existing.LastEnrollmentDate=t.LastEnrollmentDate
	existing.Price=t.Price
	existing.Discount=t.Discount
	existing.Status=t.Status
	return nil
}

func (m *MockTourRepository)DeleteTour(tourID int) error {
	if _,ok := m.tours[tourID]; !ok {
		return errors.New("tour not found")
	}
	delete(m.tours,tourID)
	return nil
}