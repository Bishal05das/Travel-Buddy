package agencyusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type DeleteAgencyUseCase struct {
	repo port.AgencyRepository
}

func NewDeleteAgencyUseCase(repo port.AgencyRepository) *DeleteAgencyUseCase {
	return &DeleteAgencyUseCase{
		repo: repo,
	}
}

func (uc *DeleteAgencyUseCase) Execute(ctx context.Context, agencyID int) error {
	return uc.repo.DeleteAgency(ctx, agencyID)
}
