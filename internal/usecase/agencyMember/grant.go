package memberusecase

import (
	"context"
	"errors"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

// ensureCanGrant stops a member from handing out permissions they do not
// hold themselves. Super users may grant anything.
func ensureCanGrant(ctx context.Context, repo port.AgencyMemberRepository, actor domain.Actor, requested []int) error {
	if actor.IsSuper() {
		return nil
	}
	held, err := repo.GetPermissionIDs(ctx, actor.ID)
	if err != nil {
		return err
	}
	allowed := make(map[int]bool, len(held))
	for _, id := range held {
		allowed[id] = true
	}
	for _, id := range requested {
		if !allowed[id] {
			return errors.New("cannot grant a permission you do not have")
		}
	}
	return nil
}
