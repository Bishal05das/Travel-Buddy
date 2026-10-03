package memberusecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
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

func (uc *CreateAgencyMemberUseCase) Execute(ctx context.Context, actor domain.Actor, req *domain.CreateMemberRequest) error {
	isOwner := strings.EqualFold(strings.TrimSpace(req.RoleName), "owner")
	if isOwner && !actor.IsSuper() {
		return domain.ErrOwnerCreationForbidden
	}
	if scope := actor.AgencyScope(); scope != nil && *scope != req.AgencyID {
		return errors.New("cannot create members for another agency")
	}
	if err := ensureCanGrant(ctx, uc.agencyMemberRepo, actor, req.Permissions); err != nil {
		return err
	}

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
		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			return errors.New("error in hashing password")
		}
		member := domain.AgencyMember{
			AgencyID: req.AgencyID,
			RoleID:   role.RoleID,
			Name:     req.Name,
			Email:    req.Email,
			Phone:    req.Phone,
			Password: hashedPassword,
			IsOwner:  isOwner,
		}
		if err := uc.agencyMemberRepo.CreateMember(txCtx, &member); err != nil {
			return fmt.Errorf("create member: %w", err)
		}
		return nil
	})
	return err

}
