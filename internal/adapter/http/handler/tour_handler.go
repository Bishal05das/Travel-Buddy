package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/bishal05das/travelbuddy/internal/validation"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

type TourHandler struct {
	createUC       port.CreateTour
	getUC          port.GetTour
	listUC         port.ListTour
	updateUC       port.UpdateTour
	updateStatusUC port.UpdateTourStatus
	deleteUC       port.DeleteTour
}

func NewTourHandler(createUC port.CreateTour, getUC port.GetTour, listUC port.ListTour, updateUC port.UpdateTour, updateStatusUC port.UpdateTourStatus, deleteUC port.DeleteTour) *TourHandler {
	return &TourHandler{
		createUC:       createUC,
		getUC:          getUC,
		listUC:         listUC,
		updateUC:       updateUC,
		updateStatusUC: updateStatusUC,
		deleteUC:       deleteUC,
	}
}

func (h *TourHandler) Create(w http.ResponseWriter, r *http.Request) {
	// var req domain.CreateTourRequest
	// err := json.NewDecoder(r.Body).Decode(&req)
	// if err != nil {
	// 	util.SendData(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// if err := validation.Validate.Struct(req); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// tour := domain.Tour{
	// 	AgencyID: req.AgencyID,
	// 	Name: req.Name,
	// 	StartDate: req.StartDate,
	// 	EndDate: req.EndDate,
	// 	AvailableSeat: req.AvailableSeat,
	// 	Description: req.Description,
	// 	LastEnrollmentDate: req.LastEnrollmentDate,
	// 	Price: req.Price,
	// 	Discount: req.Discount,
	// }
	// err = h.createUC.Execute(r.Context(),&tour)
	// if err != nil {
	// 	util.SendData(w,err.Error(),http.StatusBadRequest)
	// 	return
	// }
	// util.SendData(w, "Successfully Created Tour", http.StatusCreated)

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		util.SendData(w, "invalid multipart form data", http.StatusBadRequest)
		return
	}
	agencyIDStr := r.PathValue("agency_id")
	agencyID, err := uuid.Parse(agencyIDStr)

	if err != nil {
		util.SendData(w, "invalid agency_id", http.StatusBadRequest)
		return
	}
	startDate, err := parseDate(r.FormValue("start_date"))
	if err != nil {
		util.SendData(w, "invalid start_date, use YYYY-MM-DD or RFC3339", http.StatusBadRequest)
		return
	}

	endDate, err := parseDate(r.FormValue("end_date"))
	if err != nil {
		util.SendData(w, "invalid end_date, use YYYY-MM-DD or RFC3339", http.StatusBadRequest)
		return
	}

	lastEnrollmentDate, err := parseDate(r.FormValue("last_enrollment_date"))
	if err != nil {
		util.SendData(w, "invalid last_enrollment_date, use YYYY-MM-DD or RFC3339", http.StatusBadRequest)
		return
	}

	req := domain.CreateTourRequest{
		AgencyID:           agencyID,
		Name:               strings.TrimSpace(r.FormValue("name")),
		StartDate:          startDate,
		EndDate:            endDate,
		Description:        strings.TrimSpace(r.FormValue("description")),
		LastEnrollmentDate: lastEnrollmentDate,
	}

	// total_seat is the capacity; available_seat is still accepted from older
	// clients, since at creation every seat is available.
	seatStr := strings.TrimSpace(r.FormValue("total_seat"))
	if seatStr == "" {
		seatStr = strings.TrimSpace(r.FormValue("available_seat"))
	}
	totalSeat, err := strconv.Atoi(seatStr)
	if err != nil {
		util.SendData(w, "invalid total_seat", http.StatusBadRequest)
		return
	}
	req.TotalSeat = totalSeat

	price, err := strconv.Atoi(r.FormValue("price"))
	if err != nil {
		util.SendData(w, "invalid price", http.StatusBadRequest)
		return
	}
	req.Price = price

	discountStr := strings.TrimSpace(r.FormValue("discount"))
	if discountStr != "" {
		discount, err := strconv.Atoi(discountStr)
		if err != nil {
			util.SendData(w, "invalid discount", http.StatusBadRequest)
			return
		}
		req.Discount = discount
	}

	if err := validation.Validate.Struct(req); err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}

	fullPath, fileName, err := saveImage(r, "image", "tours")
	if err != nil {
		sendUploadError(w, err)
		return
	}

	tour := domain.Tour{
		AgencyID:           agencyID,
		Name:               req.Name,
		StartDate:          startDate,
		EndDate:            endDate,
		TotalSeat:          req.TotalSeat,
		Description:        req.Description,
		LastEnrollmentDate: lastEnrollmentDate,
		Price:              req.Price,
		Discount:           req.Discount,
		ImagePath:          filepath.ToSlash(filepath.Join("tours", fileName)),
	}

	if err := h.createUC.Execute(r.Context(), &tour); err != nil {
		_ = os.Remove(fullPath)
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.SendData(w, map[string]any{
		"message":               "Successfully Created Tour",
		"background_image_path": tour.ImagePath,
		"background_image_url":  "/images/" + tour.ImagePath,
	}, http.StatusCreated)
}

func (h *TourHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("tour_id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tour, err := h.getUC.Execute(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, tour, http.StatusOK)

}

func (h *TourHandler) List(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("agency_id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var page, limit int

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "invalid page", http.StatusBadRequest)
			return
		}
		page = p
	}
	limitStr := r.URL.Query().Get("limit")
	if limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l < 1 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = l
	}
	data, err := h.listUC.Execute(r.Context(), id, page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, data, http.StatusOK)
}

func (h *TourHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("tour_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	actor, ok := actorFromRequest(w, r)
	if !ok {
		return
	}
	err = h.deleteUC.Execute(r.Context(), actor, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, "tour Deleted Successsfully", http.StatusOK)
}

func (h *TourHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("tour_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	agencyidStr := r.PathValue("agency_id")
	agencyid, err := uuid.Parse(agencyidStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req domain.UpdateTourRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.AgencyID = agencyid
	if err := validation.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tour := &domain.Tour{
		TourID:             id,
		AgencyID:           req.AgencyID,
		Name:               req.Name,
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		TotalSeat:          req.TotalSeat,
		Description:        req.Description,
		LastEnrollmentDate: req.LastEnrollmentDate,
		Price:              req.Price,
		Discount:           req.Discount,
		UpdatedAt:          time.Now(),
	}
	err = h.updateUC.Execute(r.Context(), tour)
	if errors.Is(err, domain.ErrCapacityBelowBooked) {
		util.SendData(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, map[string]any{
		"message":        "Tour Updated Successfully",
		"total_seat":     tour.TotalSeat,
		"available_seat": tour.AvailableSeat,
	}, http.StatusOK)
}

func (h *TourHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("tour_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var status string
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	actor, ok := actorFromRequest(w, r)
	if !ok {
		return
	}
	change, err := h.updateStatusUC.Execute(r.Context(), actor, id, status)
	if errors.Is(err, domain.ErrTourCancelled) {
		util.SendData(w, err.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, change, http.StatusOK)
}

// parseDate accepts either a plain date (YYYY-MM-DD) or a full RFC3339 timestamp.
func parseDate(v string) (time.Time, error) {
	v = strings.TrimSpace(v)
	if t, err := time.Parse(time.DateOnly, v); err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339, v)
}
