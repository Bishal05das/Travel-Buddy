package tourusecase

import (
	"context"
	"sync"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

type listTourUseCase struct {
	repo port.TourRepository
}

func NewListTourUseCase(repo port.TourRepository) port.ListTour {
	return &listTourUseCase{
		repo: repo,
	}
}

func (uc *listTourUseCase) Execute(ctx context.Context, agencyID uuid.UUID, page, limit int) (*util.PaginationData, error) {
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		firstErr error
	)

	if page < 1 {
		page = 1
	}

	if limit > 100 {
		limit = 100
	}

	var list []*domain.Tour
	var count int
	capture := func(err error) {
		mu.Lock()
		defer mu.Unlock()
		if firstErr == nil {
			firstErr = err
		}
	}
	wg.Add(2)

	go func() {
		defer wg.Done()
		result, err := uc.repo.ListTour(ctx, agencyID, page, limit)
		if err != nil {
			capture(err)
			return
		}
		list = result
	}()
	go func() {
		defer wg.Done()
		cnt, err := uc.repo.Count(ctx, agencyID)
		if err != nil {
			capture(err)
			return
		}
		count = cnt
	}()
	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}
	paginationdata := util.PaginationData{
		Data: list,
		Meta: util.Meta{
			Page:       page,
			Limit:      limit,
			TotalCount: count,
			TotalPage:  (count + limit - 1) / limit,
		},
	}
	return &paginationdata, nil
}
