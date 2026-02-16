package agencyusecase

import (
	"context"
	"fmt"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type CreateAgencyUseCase struct {
	repo port.AgencyRepository
}

func NewCreateTourUseCase(r port.AgencyRepository) *CreateAgencyUseCase {
	return &CreateAgencyUseCase{
		repo: r,
	}
}

func (uc *CreateAgencyUseCase) Execute(ctx context.Context, agency *domain.Agency) error {
	//add business logic here
	if agency.Name == "" ||
		agency.Address == "" ||
		agency.RegistrationID == "" {
		return fmt.Errorf("invalid or missing agency data")
	}
	return uc.repo.CreateAgency(ctx, agency)
}
