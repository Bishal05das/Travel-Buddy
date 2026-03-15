package agencyusecase

import (
	"context"
	"errors"

	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type DeleteAgencyUseCase struct {
	repo port.AgencyRepository
}

func NewDeleteAgencyUseCase(repo port.AgencyRepository) *DeleteAgencyUseCase {
	return &DeleteAgencyUseCase{
		repo: repo,
	}
}

func (uc *DeleteAgencyUseCase) Execute(ctx context.Context, agencyID uuid.UUID) error {
	if agencyID == uuid.Nil {
		return errors.New("agencyID is required")
	}
	return uc.repo.DeleteAgency(ctx, agencyID)
}
