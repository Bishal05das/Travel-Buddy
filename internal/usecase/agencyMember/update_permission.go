package memberusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/google/uuid"
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

func (uc *UpdatPermissionUseCase) Execute(ctx context.Context, memberID uuid.UUID, req *domain.UpdatePermissionRequest) error {
	return uc.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		//apply row level locking
		roleID, err := uc.agencyMemberRepo.GetRoleIDFromMemberIDForUpdate(txCtx, memberID)
		if err != nil {
			return err
		}
		err = uc.roleRepo.DeletePermissionsFromRole(txCtx, *roleID)
		if err != nil {
			return err
		}
		// Returning this error rolls back the delete above; previously it was
		// swallowed and the member was left with no permissions at all.
		return uc.roleRepo.AddPermissionsToRole(txCtx, *roleID, req.Permissions)
	})
}
