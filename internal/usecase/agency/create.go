package agencyusecase

import (
	"context"
	"errors"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type CreateAgencyUseCase struct {
	repo port.AgencyRepository
}

func NewCreateAgencyUseCase(r port.AgencyRepository) *CreateAgencyUseCase {
	return &CreateAgencyUseCase{
		repo: r,
	}
}

func (uc *CreateAgencyUseCase) Execute(ctx context.Context, agency *domain.Agency) error {
	//add business logic here
	if agency.Name == "" ||
		agency.Address == "" ||
		agency.RegistrationID == "" {
		return errors.New("invalid or missing agency data")
	}
	return uc.repo.CreateAgency(ctx, agency)
}
