package agencyusecase

import (
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

func (uc *DeleteAgencyUseCase) Execute(agencyID int) error {
	return uc.repo.DeleteAgency(agencyID)
}
