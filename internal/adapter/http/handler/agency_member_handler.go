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
	loginUC                  port.LoginMember
}

func NewMemberHandler(createMemberUC port.CreateAgencyMember, deleteMemberUC port.DeleteAgencyMember, listMemberUC port.ListAgencyMember, updateMemberPermissionUC port.UpdateAgencyMemberPermission, loginUC port.LoginMember) *MemberHandler {
	return &MemberHandler{
		createMemberUC:           createMemberUC,
		deleteMemberUC:           deleteMemberUC,
		listMemberUC:             listMemberUC,
		updateMemberPermissionUC: updateMemberPermissionUC,
		loginUC:                  loginUC,
	}
}

func (h *MemberHandler) CreateMember(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateMemberRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.createMemberUC.Execute(r.Context(), &req)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, "Successfully Created Member", http.StatusCreated)
}

func (h *MemberHandler) DeleteMember(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("member_id")
	memberID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid agency id", http.StatusBadRequest)
		return
	}
	err = h.deleteMemberUC.Execute(r.Context(), memberID)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, "Successfully Deleted Member", http.StatusCreated)
}

func (h *MemberHandler) ListMember(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("agency_id")
	agencyID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid agency id", http.StatusBadRequest)
		return
	}
	result, err := h.listMemberUC.Execute(r.Context(), agencyID)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, result, http.StatusCreated)
}

func (h *MemberHandler) UpdateMemberPermissions(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("member_id")
	memberID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid member id", http.StatusBadRequest)
		return
	}
	var req domain.UpdatePermissionRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
	}
	err = h.updateMemberPermissionUC.Execute(r.Context(), memberID, &req)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
	}
	util.SendData(w, "Successfully Updated Permission", http.StatusCreated)
}

func (h *MemberHandler) MemberLogin(w http.ResponseWriter, r *http.Request) {
	var reqLogin domain.ReqLogin

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := h.loginUC.Execute(r.Context(), &reqLogin)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusInternalServerError)
		return
	}
	util.SendData(w, token, http.StatusCreated)
}
