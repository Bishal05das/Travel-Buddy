package mocks

import (
	"context"
	"errors"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/google/uuid"
)

type MockTourRepository struct {
	tours map[uuid.UUID]*domain.Tour
	err   error
}

func NewMockTourRepository() *MockTourRepository {
	return &MockTourRepository{
		tours: map[uuid.UUID]*domain.Tour{},
	}
}

func (m *MockTourRepository) CreateTour(ctx context.Context,tour *domain.Tour) error {
	if m.err != nil {
		return m.err
	}
	tour.TourID = uuid.New()
	m.tours[tour.TourID] = tour
	return nil
}

func (m *MockTourRepository) ListTour(ctx context.Context,agencyID uuid.UUID) ([]*domain.Tour, error) {
	var result []*domain.Tour

	for _, tour := range m.tours {
		if tour.AgencyID == agencyID {
			result = append(result, tour)
		}
	}
	return result, nil
}

func (m *MockTourRepository) UpdateTour(ctx context.Context,t *domain.Tour) error {

	existing, ok := m.tours[t.TourID]
	if !ok {
		return errors.New("tour not found")
	}
	existing.AgencyID = t.AgencyID
	existing.Name = t.Name
	existing.StartDate = t.StartDate
	existing.EndDate = t.EndDate
	existing.Description = t.Description
	existing.LastEnrollmentDate = t.LastEnrollmentDate
	existing.Price = t.Price
	existing.Discount = t.Discount
	existing.Status = t.Status
	return nil
}

func (m *MockTourRepository) DeleteTour(ctx context.Context,tourID uuid.UUID) error {
	if _, ok := m.tours[tourID]; !ok {
		return errors.New("tour not found")
	}
	delete(m.tours, tourID)
	return nil
}

func (m *MockTourRepository) GetByID(ctx context.Context, tourID uuid.UUID) (*domain.Tour, error)
func (m *MockTourRepository) GetByIDForUpdate(ctx context.Context, tourID uuid.UUID) (*domain.Tour, error)

func (m *MockTourRepository) UpdateAvailableSeats(ctx context.Context, tourID uuid.UUID, seats int) error
