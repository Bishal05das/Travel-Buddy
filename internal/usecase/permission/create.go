package permissionusecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type createPermissionUseCase struct {
	repo port.PermissionRepository
}

func NewCreatePermissionUseCase(repo port.PermissionRepository) port.CreatePermission {
	return &createPermissionUseCase{repo: repo}
}

func (uc *createPermissionUseCase) Execute(ctx context.Context, p *domain.Permission) error {
	if p == nil {
		return errors.New("permission payload is required")
	}

	if p.Name == "" || p.Resource == "" || p.Action == "" {
		return errors.New("missing required permission fields")
	}

	if err := uc.repo.CreatePermission(ctx, p); err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}

	return nil
}
