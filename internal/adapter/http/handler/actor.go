package handler

import (
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/domain"
	util "github.com/bishal05das/travelbuddy/utils"
)

// actorFromRequest returns the authenticated caller. ok is false (and a 401
// has been written) when the request carries no verified claims.
func actorFromRequest(w http.ResponseWriter, r *http.Request) (domain.Actor, bool) {
	p, err := util.GetPayload(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return domain.Actor{}, false
	}
	return domain.Actor{ID: p.UserID, Role: p.Role, AgencyID: p.AgencyID}, true
}
