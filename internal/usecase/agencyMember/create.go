package memberusecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type CreateAgencyMemberUseCase struct {
	txManager        port.TxManager
	agencyMemberRepo port.AgencyMemberRepository
	roleRepo         port.RoleRepository
}

func NewCreateAgencyMemberUseCase(txManager port.TxManager, agencyMemberRepo port.AgencyMemberRepository, roleRepo port.RoleRepository) *CreateAgencyMemberUseCase {
	return &CreateAgencyMemberUseCase{
		txManager:        txManager,
		agencyMemberRepo: agencyMemberRepo,
		roleRepo:         roleRepo,
	}
}

func (uc *CreateAgencyMemberUseCase) Execute(ctx context.Context, req *domain.CreateMemberRequest) error {

	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		//create role
		role := domain.Role{
			AgencyID: req.AgencyID,
			RoleName: req.RoleName,
		}
		err := uc.roleRepo.CreateRole(txCtx, &role)
		if err != nil {
			return errors.New(err.Error())
		}
		//add permissions to role
		err = uc.roleRepo.AddPermissionsToRole(txCtx, role.RoleID, req.Permissions)
		if err != nil { 
			return errors.New(err.Error())
		}
		member := domain.AgencyMember{
			AgencyID: req.AgencyID,
			RoleID:   role.RoleID,
			Name:     req.Name,
			Email:    req.Email,
			Phone:    req.Phone,
			Password: req.Password,
		}
		if err := uc.agencyMemberRepo.CreateMember(txCtx, &member); err != nil {
			return fmt.Errorf("create member: %w", err)
		}
		return nil
	})
	return err

}
