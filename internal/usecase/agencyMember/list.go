package agencymember

import (
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type ListAgencyMemberUseCase struct {
	repo port.AgencyMemberRepository
}

func NewListAgencyMemberUseCase(repo port.AgencyMemberRepository) *ListAgencyMemberUseCase {
	return &ListAgencyMemberUseCase{
		repo: repo,
	}
}

func(uc *ListAgencyMemberUseCase) Execute(agencyID int) ([]*domain.AgencyMember,error) {
	return uc.repo.ListMember(agencyID)
}