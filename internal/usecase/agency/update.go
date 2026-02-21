package agencyusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type UpdateAgencyUseCase struct {
	repo port.AgencyRepository
}

func NewUpdateAgencyUseCase(repo port.AgencyRepository) *UpdateAgencyUseCase {
	return &UpdateAgencyUseCase{
		repo: repo,
	}
}

func (uc *UpdateAgencyUseCase) Execute(ctx context.Context, agency *domain.Agency) error {

	return uc.repo.UpdateAgency(ctx, agency)
}
