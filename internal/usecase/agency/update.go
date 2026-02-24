package agencyusecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
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
	if agency == nil {
		return errors.New("agency payload is required")
	}

	if agency.AgencyID == uuid.Nil {
		return errors.New("Invalid agency id")
	}

	if agency.Name == "" ||
		agency.Address == "" ||
		agency.RegistrationID == "" {
		return errors.New("missing required agency fields")
	}

	// Check if agency exists
	// existing, err := uc.repo.GetAgencyByID(ctx, agency.AgencyID)
	// if err != nil {
	// 	return fmt.Errorf("failed to fetch agency: %w", err)
	// }
	// if existing == nil {
	// 	return errors.New("agency not found")
	// }

	// Business rule example
	// Prevent updating inactive agencies (optional rule)
	// if !existing.IsActive {
	// 	return errors.New("cannot update inactive agency")
	// }

	// Preserve immutable fields (if any)
	// agency.IsActive = existing.IsActive

	// Perform update
	if err := uc.repo.UpdateAgency(ctx, agency); err != nil {
		return fmt.Errorf("failed to update agency: %w", err)
	}

	return nil
}
