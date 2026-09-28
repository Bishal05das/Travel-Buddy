package memberusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type DeleteAgencyMemberUseCase struct {
	repo port.AgencyMemberRepository
}

func NewDeleteAgencyMemberUseCase(repo port.AgencyMemberRepository) *DeleteAgencyMemberUseCase {
	return &DeleteAgencyMemberUseCase{
		repo: repo,
	}
}

func (uc *DeleteAgencyMemberUseCase) Execute(ctx context.Context, actor domain.Actor, agencyMemberID uuid.UUID) error {
	return uc.repo.DeleteMember(ctx, agencyMemberID, actor.AgencyScope())
}
