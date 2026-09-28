package permissionusecase

import (
	"context"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type listPermissionsUseCase struct {
	repo port.PermissionRepository
}

func NewListPermissionsUseCase(repo port.PermissionRepository) port.ListPermissions {
	return &listPermissionsUseCase{repo: repo}
}

func (uc *listPermissionsUseCase) Execute(ctx context.Context) ([]domain.Permission, error) {
	return uc.repo.ListPermissions(ctx)
}
