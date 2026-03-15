package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
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
	var tour domain.Tour
	err := json.NewDecoder(r.Body).Decode(&tour)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	err = h.createUC.Execute(r.Context(),&tour)
	if err != nil {
		util.SendData(w,err.Error(),http.StatusBadRequest)
		return
	}
	util.SendData(w, tour, http.StatusCreated)
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
	var tour domain.Tour
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&tour)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	tour.TourID=id
	err = h.updateUC.Execute(r.Context(), &tour)
	fmt.Println(err)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, "Tour Updated Successfully", http.StatusOK)
}