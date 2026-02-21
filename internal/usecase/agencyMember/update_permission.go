package memberusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type UpdatPermissionUseCase struct {
	txManager        port.TxManager
	agencyMemberRepo port.AgencyMemberRepository
	roleRepo         port.RoleRepository
}

func NewUpdatePermissionUseCase(txManager port.TxManager, agencyMemberRepo port.AgencyMemberRepository, roleRepo port.RoleRepository) *UpdatPermissionUseCase {
	return &UpdatPermissionUseCase{
		txManager:        txManager,
		agencyMemberRepo: agencyMemberRepo,
		roleRepo:         roleRepo,
	}
}

func (uc *UpdatPermissionUseCase) Execute(ctx context.Context, req *domain.UpdatePermissionRequest) error {
	var response error
	err := uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		roleID, err := uc.agencyMemberRepo.GetRoleIDFromMemberID(txCtx, req.MemberID)
		if err != nil {
			return err
		}
		err = uc.roleRepo.DeletePermissionsFromRole(txCtx, *roleID)
		if err != nil {
			return err
		}
		response = uc.roleRepo.AddPermissionsToRole(txCtx, *roleID, req.Permissions)
		return err
	})
	if err != nil {
		return err
	}
	return response
}
