package domain

import "github.com/google/uuid"

const (
	RoleSuper  = "super"  // platform administrator
	RoleUser   = "user"   // customer account
	RoleMember = "member" // agency member, limited to their own agency
)

// Actor is the authenticated caller of a use case.
type Actor struct {
	ID       uuid.UUID
	Role     string
	AgencyID *uuid.UUID // set for agency members
}

func (a Actor) IsSuper() bool { return a.Role == RoleSuper }

// AgencyScope restricts agency-owned records to the actor's agency. It is
// nil for super users, who may act on any agency.
func (a Actor) AgencyScope() *uuid.UUID {
	if a.IsSuper() {
		return nil
	}
	if a.AgencyID == nil {
		// Never nil for a non-super actor: an unknown agency matches nothing.
		none := uuid.Nil
		return &none
	}
	return a.AgencyID
}
