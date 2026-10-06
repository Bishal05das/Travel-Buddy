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
	if err := filter.Prepare(); err != nil {
		return nil, err
	}

	tours, err := uc.repo.SearchTours(ctx, filter)
	if err != nil {
		return nil, err
	}
	count, err := uc.repo.CountTours(ctx, filter)
	if err != nil {
		return nil, err
	}

	agencies, err := uc.repo.SearchAgencies(ctx, filter.Query, 5)
	if err != nil {
		return nil, err
	}
	totalPages := count / filter.Limit
	if count%filter.Limit != 0 {
		totalPages++
	}

	return &domain.SearchResult{
		Tours:    tours,
		Agencies: agencies,
		Meta:     domain.SearchMeta{Page: filter.Page, Limit: filter.Limit, TotalCount: count, TotalPage: totalPages},
	}, nil
}
