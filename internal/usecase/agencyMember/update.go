package agencymember

import (
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type UpdateAgencyMemberUseCase struct {
	repo port.AgencyMemberRepository
}

func NewUpdateAgencyMemberUseCase(repo port.AgencyMemberRepository) *UpdateAgencyMemberUseCase {
	return &UpdateAgencyMemberUseCase{
		repo: repo,
	}
}

func (uc *UpdateAgencyMemberUseCase)Execute(member *domain.AgencyMember) error {

	return uc.repo.UpdateMember(member)
}