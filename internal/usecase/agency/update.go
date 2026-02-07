package agencyusecase

import (
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type UpdateAgencyUseCase struct {
	repo port.AgencyRepository
}

func NewUpdateTourUseCase(repo port.AgencyRepository) *UpdateAgencyUseCase {
	return &UpdateAgencyUseCase{
		repo: repo,
	}
}

func (uc *UpdateAgencyUseCase)Execute(agency *domain.Agency) error {

	return uc.repo.UpdateAgency(agency)
}