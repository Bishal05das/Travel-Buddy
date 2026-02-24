package permissionusecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type deletePermissionUseCase struct {
	repo port.PermissionRepository
}

func NewDeletePermissionUseCase(repo port.PermissionRepository) port.DeletePermission {
	return &deletePermissionUseCase{repo: repo}
}

func (uc *deletePermissionUseCase) Execute(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid permission id")
	}

	if err := uc.repo.DeletePermission(ctx, id); err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	return nil
}
