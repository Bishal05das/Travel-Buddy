package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/bishal05das/travelbuddy/internal/validation"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

type MemberHandler struct {
	createMemberUC           port.CreateAgencyMember
	deleteMemberUC           port.DeleteAgencyMember
	listMemberUC             port.ListAgencyMember
	updateMemberPermissionUC port.UpdateAgencyMemberPermission
	loginUC                  port.LoginMember
	profileUC                port.GetAgencyMemberProfile
}

func NewMemberHandler(createMemberUC port.CreateAgencyMember, deleteMemberUC port.DeleteAgencyMember, listMemberUC port.ListAgencyMember, updateMemberPermissionUC port.UpdateAgencyMemberPermission, loginUC port.LoginMember, profileUC port.GetAgencyMemberProfile) *MemberHandler {
	return &MemberHandler{
		createMemberUC:           createMemberUC,
		deleteMemberUC:           deleteMemberUC,
		listMemberUC:             listMemberUC,
		updateMemberPermissionUC: updateMemberPermissionUC,
		loginUC:                  loginUC,
		profileUC:                profileUC,
	}
}

func (h *MemberHandler) CreateMember(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("agency_id")
	agencyID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid agency id", http.StatusBadRequest)
		return
	}
	var req domain.CreateMemberRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.AgencyID = agencyID
	if err := validation.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actor, ok := actorFromRequest(w, r)
	if !ok {
		return
	}
	err = h.createMemberUC.Execute(r.Context(), actor, &req)
	if err != nil {
		sendMemberError(w, err)
		return
	}
	util.SendData(w, "Successfully Created Member", http.StatusCreated)
}

func (h *MemberHandler) DeleteMember(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("member_id")
	memberID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid member id", http.StatusBadRequest)
		return
	}
	actor, ok := actorFromRequest(w, r)
	if !ok {
		return
	}
	err = h.deleteMemberUC.Execute(r.Context(), actor, memberID)
	if err != nil {
		sendMemberError(w, err)
		return
	}
	util.SendData(w, "Successfully Deleted Member", http.StatusOK)
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
	util.SendData(w, result, http.StatusOK)
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
		return
	}
	if err := validation.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	actor, ok := actorFromRequest(w, r)
	if !ok {
		return
	}
	err = h.updateMemberPermissionUC.Execute(r.Context(), actor, memberID, &req)
	if err != nil {
		sendMemberError(w, err)
		return
	}
	util.SendData(w, "Successfully Updated Permission", http.StatusOK)
}

func sendMemberError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if errors.Is(err, domain.ErrOwnerProtected) || errors.Is(err, domain.ErrOwnerCreationForbidden) {
		status = http.StatusForbidden
	}
	util.SendData(w, err.Error(), status)
}

// GetMyProfile uses verified claims, never a member ID supplied by the client.
func (h *MemberHandler) GetMyProfile(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorFromRequest(w, r)
	if !ok {
		return
	}
	if actor.Role != domain.RoleMember || actor.AgencyID == nil {
		util.SendData(w, "agency member required", http.StatusForbidden)
		return
	}
	profile, err := h.profileUC.Execute(r.Context(), actor)
	if errors.Is(err, domain.ErrMemberNotFound) {
		util.SendData(w, "member account no longer exists", http.StatusUnauthorized)
		return
	}
	if err != nil {
		util.SendData(w, "could not load profile", http.StatusInternalServerError)
		return
	}
	util.SendData(w, profile, http.StatusOK)
}

func (h *MemberHandler) MemberLogin(w http.ResponseWriter, r *http.Request) {
	var req domain.ReqLogin

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validation.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := h.loginUC.Execute(r.Context(), &req)
	if err != nil {
		sendLoginError(w, err)
		return
	}
	util.SendData(w, token, http.StatusOK)
}
