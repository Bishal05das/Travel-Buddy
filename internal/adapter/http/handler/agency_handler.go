package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/bishal05das/travelbuddy/internal/validation"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

type AgencyHandler struct {
	createUC port.CreateAgency
	updateUC port.UpdateAgency
	deleteUC port.DeleteAgency
	getUC    port.GetAgency
}

func NewAgencyHandler(createUC port.CreateAgency, updateUC port.UpdateAgency, deleteUC port.DeleteAgency, getUC port.GetAgency) *AgencyHandler {
	return &AgencyHandler{
		createUC: createUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
		getUC:    getUC,
	}
}

// GetAgency returns an agency's public profile.
func (h *AgencyHandler) GetAgency(w http.ResponseWriter, r *http.Request) {
	agencyID, err := uuid.Parse(r.PathValue("agency_id"))
	if err != nil {
		http.Error(w, "invalid agency id", http.StatusBadRequest)
		return
	}
	agency, err := h.getUC.Execute(r.Context(), agencyID)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusNotFound)
		return
	}
	util.SendData(w, agency, http.StatusOK)
}

func (h *AgencyHandler) CreateAgency(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		util.SendData(w, "invalid multipart form data", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	address := strings.TrimSpace(r.FormValue("address"))
	registrationID := strings.TrimSpace(r.FormValue("registration_id"))

	if name == "" || address == "" || registrationID == "" {
		util.SendData(w, "name, address, and registration_id are required", http.StatusBadRequest)
		return
	}

	dstPath, _, err := saveImage(r, "image", "agencies")
	if err != nil {
		sendUploadError(w, err)
		return
	}

	agency := &domain.Agency{
		Name:           name,
		Address:        address,
		RegistrationID: registrationID,
	}
	imagePath := filepath.ToSlash(dstPath) // images/agencies/abc.png

	if err := h.createUC.Execute(r.Context(), agency, imagePath); err != nil {
		_ = os.Remove(dstPath) // rollback file if DB save fails
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.SendData(w, map[string]any{
		"message":    "Agency successfully created",
		"agency_id":  agency.AgencyID,
		"image_path": imagePath,
		"image_url":  "/" + imagePath,
	}, http.StatusCreated)

}

func (h *AgencyHandler) UpdateAgency(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("agency_id")
	AgencyID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid agency id", http.StatusBadRequest)
		return
	}
	var req domain.UpdateAgencyRequest
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
	agency := &domain.Agency{
		AgencyID:       AgencyID,
		Name:           req.Name,
		Address:        req.Address,
		RegistrationID: req.RegistrationID,
		UpdatedAt:      time.Now(),
	}
	err = h.updateUC.Execute(r.Context(), agency)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, "Agency Updated Successfully", http.StatusOK)

}

func (h *AgencyHandler) DeleteAgency(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("agency_id")
	agencyID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid agency id", http.StatusBadRequest)
		return
	}
	err = h.deleteUC.Execute(r.Context(), agencyID)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, "Agency Deleted Successfully", http.StatusOK)

}
