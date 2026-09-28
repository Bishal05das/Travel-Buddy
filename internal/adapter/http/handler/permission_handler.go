package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
)

type PermissionHandler struct {
	createUC port.CreatePermission
	deleteUC port.DeletePermission
	listUC   port.ListPermissions
}

func NewPermissionHandler(c port.CreatePermission, d port.DeletePermission, l port.ListPermissions) *PermissionHandler {
	return &PermissionHandler{
		createUC: c,
		deleteUC: d,
		listUC:   l,
	}
}

// ListPermissions returns the permission catalogue, so clients can offer
// permissions by name when creating members.
func (h *PermissionHandler) ListPermissions(w http.ResponseWriter, r *http.Request) {
	perms, err := h.listUC.Execute(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	util.SendData(w, perms, http.StatusOK)
}

func (h *PermissionHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var perm domain.Permission

	if err := json.NewDecoder(r.Body).Decode(&perm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.createUC.Execute(r.Context(), &perm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.SendData(w, perm, http.StatusCreated)
}

func (h *PermissionHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid permission id", http.StatusBadRequest)
		return
	}

	if err := h.deleteUC.Execute(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.SendData(w, "Permission Deleted Succesfully", http.StatusOK)
}
