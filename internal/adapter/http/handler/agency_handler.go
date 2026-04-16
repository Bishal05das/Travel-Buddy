package handler

import (
	"encoding/json"
	"fmt"
	"io"
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
}

func NewAgencyHandler(createUC port.CreateAgency, updateUC port.UpdateAgency, deleteUC port.DeleteAgency) *AgencyHandler {
	return &AgencyHandler{
		createUC: createUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
	}
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

	file, header, err := r.FormFile("image")
	if err != nil {
		util.SendData(w, "agency image is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowed[ext] {
		util.SendData(w, "only jpg, jpeg, png, and webp files are allowed", http.StatusBadRequest)
		return
	}

	uploadDir := filepath.Join("images", "agencies")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		fmt.Println(err)
		util.SendData(w, "failed to create upload directory", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.NewString(), ext)
	dstPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		util.SendData(w, "failed to save image", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		util.SendData(w, "failed to write image", http.StatusInternalServerError)
		return
	}

	agency := &domain.Agency{
		Name:           name,
		Address:        address,
		RegistrationID: registrationID,
		ImagePath:      filepath.ToSlash(dstPath), // images/agencies/abc.png
	}

	if err := h.createUC.Execute(r.Context(), agency); err != nil {
		_ = os.Remove(dstPath) // rollback file if DB save fails
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.SendData(w, map[string]any{
		"message":    "Agency successfully created",
		"image_path": agency.ImagePath,
		"image_url":  "/" + filepath.ToSlash(dstPath),
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
