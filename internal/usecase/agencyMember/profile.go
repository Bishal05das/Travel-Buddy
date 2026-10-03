package memberusecase

import (
	"context"
	"errors"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type GetAgencyMemberProfileUseCase struct {
	repo port.AgencyMemberRepository
}

func NewGetAgencyMemberProfileUseCase(repo port.AgencyMemberRepository) *GetAgencyMemberProfileUseCase {
	return &GetAgencyMemberProfileUseCase{repo: repo}
}

func (uc *GetAgencyMemberProfileUseCase) Execute(ctx context.Context, actor domain.Actor) (*domain.MemberProfile, error) {
	if actor.Role != domain.RoleMember || actor.AgencyID == nil {
		return nil, errors.New("agency member required")
	}
	return uc.repo.GetMemberProfile(ctx, actor.ID, *actor.AgencyID)
}
