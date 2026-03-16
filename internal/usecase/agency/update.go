package agencyusecase

import (
	"context"
	"fmt"

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

	if err := uc.repo.UpdateAgency(ctx, agency); err != nil {
		return fmt.Errorf("failed to update agency: %w", err)
	}

	return nil
}
