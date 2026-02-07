package agencymember

import (
	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type CreateAgencyMemberUseCase struct {
	repo port.AgencyMemberRepository
}

func NewCreateAgencyMemberUseCase(repo port.AgencyMemberRepository) *CreateAgencyMemberUseCase{
	return &CreateAgencyMemberUseCase{
		repo: repo,
	}
}

func(uc *CreateAgencyMemberUseCase)Execute(member *domain.AgencyMember) error {
	return uc.repo.CreateMember(member)
}