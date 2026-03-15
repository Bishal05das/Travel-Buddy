package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
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
	var agency domain.Agency
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&agency)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.createUC.Execute(r.Context(), &agency)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, "Agent Successfully Created", http.StatusCreated)

}

func (h *AgencyHandler) UpdateAgency(w http.ResponseWriter, r *http.Request) {
	var agency domain.Agency
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&agency)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.updateUC.Execute(r.Context(), &agency)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, "Agent Updated Successfully", http.StatusOK)

}

func (h *AgencyHandler) DeleteAgency(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
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
	util.SendData(w, "Agent Deleted Successfully", http.StatusOK)

}
