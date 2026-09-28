package mocks

import (
	"context"
	"errors"
	"sort"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

var _ port.TourRepository = (*MockTourRepository)(nil)

var errTourNotFound = errors.New("tour not found")

// MockTourRepository is an in-memory port.TourRepository for use case tests.
type MockTourRepository struct {
	tours map[uuid.UUID]*domain.Tour
	err   error
}

func NewMockTourRepository() *MockTourRepository {
	return &MockTourRepository{
		tours: map[uuid.UUID]*domain.Tour{},
	}
}

// SetError makes every subsequent call fail with err.
func (m *MockTourRepository) SetError(err error) {
	m.err = err
}

func (m *MockTourRepository) CreateTour(ctx context.Context, tour *domain.Tour) error {
	if m.err != nil {
		return m.err
	}
	// Keep a caller-supplied ID so tests can seed tours they refer to later.
	if tour.TourID == uuid.Nil {
		tour.TourID = uuid.New()
	}
	if tour.Status == "" {
		tour.Status = "open"
	}
	m.tours[tour.TourID] = tour
	return nil
}

func (m *MockTourRepository) ListTour(ctx context.Context, agencyID uuid.UUID, page, limit int) ([]*domain.Tour, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []*domain.Tour
	for _, tour := range m.tours {
		if tour.AgencyID == agencyID {
			result = append(result, tour)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].StartDate.After(result[j].StartDate) })

	offset := (page - 1) * limit
	if offset >= len(result) {
		return []*domain.Tour{}, nil
	}
	end := offset + limit
	if end > len(result) {
		end = len(result)
	}
	return result[offset:end], nil
}

func (m *MockTourRepository) Count(ctx context.Context, agencyID uuid.UUID) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	count := 0
	for _, tour := range m.tours {
		if tour.AgencyID == agencyID {
			count++
		}
	}
	return count, nil
}

func (m *MockTourRepository) UpdateTour(ctx context.Context, t *domain.Tour) error {
	if m.err != nil {
		return m.err
	}
	existing, ok := m.tours[t.TourID]
	if !ok || existing.AgencyID != t.AgencyID {
		return errTourNotFound
	}
	existing.Name = t.Name
	existing.StartDate = t.StartDate
	existing.EndDate = t.EndDate
	existing.AvailableSeat = t.AvailableSeat
	existing.Description = t.Description
	existing.LastEnrollmentDate = t.LastEnrollmentDate
	existing.Price = t.Price
	existing.Discount = t.Discount
	existing.UpdatedAt = t.UpdatedAt
	return nil
}

func (m *MockTourRepository) DeleteTour(ctx context.Context, tourID uuid.UUID, agencyScope *uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	if tour, ok := m.tours[tourID]; !ok || !inScope(tour, agencyScope) {
		return errTourNotFound
	}
	delete(m.tours, tourID)
	return nil
}

func (m *MockTourRepository) GetByID(ctx context.Context, tourID uuid.UUID) (*domain.Tour, error) {
	if m.err != nil {
		return nil, m.err
	}
	tour, ok := m.tours[tourID]
	if !ok {
		return nil, errTourNotFound
	}
	copied := *tour
	return &copied, nil
}

func (m *MockTourRepository) GetByIDForUpdate(ctx context.Context, tourID uuid.UUID) (*domain.Tour, error) {
	return m.GetByID(ctx, tourID)
}

func (m *MockTourRepository) UpdateAvailableSeats(ctx context.Context, tourID uuid.UUID, seats int) error {
	if m.err != nil {
		return m.err
	}
	tour, ok := m.tours[tourID]
	if !ok {
		return errTourNotFound
	}
	tour.AvailableSeat = seats
	return nil
}

func (m *MockTourRepository) UpdateTourStatus(ctx context.Context, tourID uuid.UUID, status string, agencyScope *uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	tour, ok := m.tours[tourID]
	if !ok || !inScope(tour, agencyScope) {
		return errTourNotFound
	}
	tour.Status = status
	return nil
}

func inScope(tour *domain.Tour, agencyScope *uuid.UUID) bool {
	return agencyScope == nil || tour.AgencyID == *agencyScope
}
