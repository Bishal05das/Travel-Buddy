package searchusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type searchUseCase struct {
	repo port.SearchRepository
}

func NewSearchUseCase(repo port.SearchRepository) *searchUseCase {
	return &searchUseCase{repo: repo}
}

func (uc *searchUseCase) Execute(ctx context.Context, filter domain.TourSearchFilter) (*domain.SearchResult, error) {

	tours, err := uc.repo.SearchTours(ctx, filter)
	if err != nil {
		return nil, err
	}

	agencies, err := uc.repo.SearchAgencies(ctx, filter.Query, 5)
	if err != nil {
		return nil, err
	}

	return &domain.SearchResult{
		Tours:    tours,
		Agencies: agencies,
	}, nil
}
