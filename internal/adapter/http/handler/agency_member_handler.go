package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

type MemberHandler struct {
	createMemberUC           port.CreateAgencyMember
	deleteMemberUC           port.DeleteAgencyMember
	listMemberUC             port.ListAgencyMember
	updateMemberPermissionUC port.UpdateAgencyMemberPermission
}

func NewMemberHandler(createMemberUC port.CreateAgencyMember, deleteMemberUC port.DeleteAgencyMember, listMemberUC port.ListAgencyMember, updateMemberPermissionUC port.UpdateAgencyMemberPermission) *MemberHandler {
	return &MemberHandler{
		createMemberUC:           createMemberUC,
		deleteMemberUC:           deleteMemberUC,
		listMemberUC:             listMemberUC,
		updateMemberPermissionUC: updateMemberPermissionUC,
	}
}

func (h *MemberHandler) CreateMember(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateMemberRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
	}
	err = h.createMemberUC.Execute(r.Context(), &req)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
	}
	util.SendData(w, "Successfully Created Member", http.StatusCreated)
}

func (h *MemberHandler) DeleteMember(w http.ResponseWriter, r *http.Request) {
	var memberID uuid.UUID
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&memberID)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
	}
	err = h.deleteMemberUC.Execute(r.Context(), memberID)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
	}
	util.SendData(w, "Successfully Deleted Member", http.StatusCreated)
}

func (h *MemberHandler) ListMember(w http.ResponseWriter, r *http.Request) {
	var agencyID uuid.UUID
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&agencyID)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
	}
	result, err := h.listMemberUC.Execute(r.Context(), agencyID)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
	}
	util.SendData(w, result, http.StatusCreated)
}

func (h *MemberHandler) UpdateMemberPermissions(w http.ResponseWriter, r *http.Request) {
	var req domain.UpdatePermissionRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
	}
	err = h.updateMemberPermissionUC.Execute(r.Context(), &req)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
	}
	util.SendData(w, "Successfully Updated Permission", http.StatusCreated)
}
