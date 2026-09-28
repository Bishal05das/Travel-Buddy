package handler_test

import (
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/domain"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

// asSuper attaches verified super-user claims, as the Authentication
// middleware would.
func asSuper(req *http.Request) *http.Request {
	return req.WithContext(util.WithPayload(req.Context(), &util.Payload{UserID: uuid.New(), Role: domain.RoleSuper}))
}
