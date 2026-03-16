package memberusecase

import (
	"context"
	"errors"

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

func (uc *ListAgencyMemberUseCase) Execute(ctx context.Context, agencyID uuid.UUID) ([]*domain.ListMemberResponse, error) {
	if agencyID == uuid.Nil {
		return nil, errors.New("AgencyID is Required")
	}
	return uc.repo.ListMember(ctx, agencyID)
}
