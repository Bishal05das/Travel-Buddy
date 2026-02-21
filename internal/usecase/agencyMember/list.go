package memberusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
)

type ListAgencyMemberUseCase struct {
	repo port.AgencyMemberRepository
}

func NewListAgencyMemberUseCase(repo port.AgencyMemberRepository) *ListAgencyMemberUseCase {
	return &ListAgencyMemberUseCase{
		repo: repo,
	}
}

func (uc *ListAgencyMemberUseCase) Execute(ctx context.Context, agencyID uuid.UUID) ([]*domain.AgencyMember, error) {
	return uc.repo.ListMember(ctx, agencyID)
}
