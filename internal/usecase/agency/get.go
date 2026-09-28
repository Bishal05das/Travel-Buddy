package agencyusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type getAgencyUseCase struct {
	repo port.AgencyRepository
}

func NewGetAgencyUseCase(repo port.AgencyRepository) port.GetAgency {
	return &getAgencyUseCase{repo: repo}
}

func (uc *getAgencyUseCase) Execute(ctx context.Context, agencyID uuid.UUID) (*domain.Agency, error) {
	return uc.repo.GetAgency(ctx, agencyID)
}
