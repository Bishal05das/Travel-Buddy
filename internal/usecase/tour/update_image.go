package tourusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type updateTourImageUseCase struct{ repo port.TourRepository }

func NewUpdateTourImageUseCase(repo port.TourRepository) port.UpdateTourImage {
	return &updateTourImageUseCase{repo: repo}
}

func (uc *updateTourImageUseCase) Execute(ctx context.Context, actor domain.Actor, agencyID, tourID uuid.UUID, imagePath string) (string, error) {
	if !actor.IsSuper() && (actor.Role != domain.RoleMember || actor.AgencyID == nil || *actor.AgencyID != agencyID) {
		return "", domain.ErrImageAccessDenied
	}
	return uc.repo.UpdateTourImage(ctx, agencyID, tourID, imagePath)
}
