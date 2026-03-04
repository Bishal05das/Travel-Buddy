package homeusecase

import (
	"context"
	"sync"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type homeUseCase struct {
	repo port.HomeRepository
}

func NewHomeUseCase(repo port.HomeRepository) port.Home {
	return &homeUseCase{
		repo: repo,
	}
}

func (u *homeUseCase) GetHome(ctx context.Context) (*domain.HomeResponse, error) {
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		firstErr error
		response domain.HomeResponse
	)

	capture := func(err error) {
		mu.Lock()
		defer mu.Unlock()
		if firstErr == nil {
			firstErr = err
		}
	}
	wg.Add(2)

	go func(){
		defer wg.Done()
		result,err := u.repo.GetTopTours(ctx,5)
		if err != nil {
			capture(err)
			return
		}
		mu.Lock()
		response.TopTours = result
		mu.Unlock()

	}()

	go func() {
		defer wg.Done()
		result, err := u.repo.GetTopAgencies(ctx, 6)
		if err != nil {
			capture(err)
			return
		}
		mu.Lock()
		response.TopAgencies = result
		mu.Unlock()
	}()

	wg.Wait()

	if firstErr != nil {
		return nil,firstErr
	}
	return &response, nil
}
