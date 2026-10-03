package agencyusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type updateAgencyImageUseCase struct{ repo port.AgencyRepository }

func NewUpdateAgencyImageUseCase(repo port.AgencyRepository) port.UpdateAgencyImage {
	return &updateAgencyImageUseCase{repo: repo}
}

func (uc *updateAgencyImageUseCase) Execute(ctx context.Context, actor domain.Actor, agencyID uuid.UUID, imagePath string) (string, error) {
	if !actor.IsSuper() && (actor.Role != domain.RoleMember || actor.AgencyID == nil || *actor.AgencyID != agencyID) {
		return "", domain.ErrImageAccessDenied
	}
	return uc.repo.UpdateAgencyImage(ctx, agencyID, imagePath)
}
