package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/bishal05das/travelbuddy/internal/validation"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

type TourHandler struct {
	createUC port.CreateTour
	getUC port.GetTour
	listUC port.ListTour
	updateUC port.UpdateTour
	deleteUC port.DeleteTour
}

func NewTourHandler(createUC port.CreateTour,getUC port.GetTour,listUC port.ListTour,updateUC port.UpdateTour,deleteUC port.DeleteTour) *TourHandler {
	return &TourHandler{
		createUC: createUC,
		getUC: getUC,
		listUC: listUC,
		updateUC: updateUC,
		deleteUC: deleteUC,
	}
}

func (h *TourHandler) Create(w http.ResponseWriter, r *http.Request)  {
	var req domain.CreateTourRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validation.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tour := domain.Tour{
		AgencyID: req.AgencyID,
		Name: req.Name,
		StartDate: req.StartDate,
		EndDate: req.EndDate,
		AvailableSeat: req.AvailableSeat,
		Description: req.Description,
		LastEnrollmentDate: req.LastEnrollmentDate,
		Price: req.Price,
		Discount: req.Discount,
	}
	err = h.createUC.Execute(r.Context(),&tour)
	if err != nil {
		util.SendData(w,err.Error(),http.StatusBadRequest)
		return
	}
	util.SendData(w, "Successfully Created Tour", http.StatusCreated)
}
func (h *TourHandler) Get(w http.ResponseWriter, r *http.Request){
	idStr :=r.PathValue("tour_id")
	id,err :=uuid.Parse(idStr)

	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	tour,err :=h.getUC.Execute(r.Context(),id)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	util.SendData(w,tour,http.StatusOK)

}

func (h *TourHandler) List(w http.ResponseWriter, r *http.Request){
	idStr :=r.PathValue("agency_id")
	id,err :=uuid.Parse(idStr)

	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	tours,err := h.listUC.Execute(r.Context(),id)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	util.SendData(w,tours,http.StatusOK)
}

func (h *TourHandler) Delete(w http.ResponseWriter, r *http.Request){
	idStr :=r.PathValue("tour_id")
	id,err :=uuid.Parse(idStr)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	err = h.deleteUC.Execute(r.Context(),id)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	util.SendData(w,"tour Deleted Successsfully",http.StatusOK)
}

func (h *TourHandler) Update(w http.ResponseWriter, r *http.Request){
	idStr :=r.PathValue("tour_id")
	id,err :=uuid.Parse(idStr)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	var req domain.UpdateTourRequest
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
	tour := &domain.Tour{
		TourID: id,
		AgencyID: req.AgencyID,
		Name: req.Name,
		StartDate: req.StartDate,
		EndDate: req.EndDate,
		AvailableSeat: req.AvailableSeat,
		Description: req.Description,
		LastEnrollmentDate: req.LastEnrollmentDate,
		Price: req.Price,
		Discount: req.Discount,
		UpdatedAt: time.Now(),
	}
	err = h.updateUC.Execute(r.Context(), tour)
	fmt.Println(err)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, "Tour Updated Successfully", http.StatusOK)
}