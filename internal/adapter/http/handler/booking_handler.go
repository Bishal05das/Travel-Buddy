package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

type BookingHandler struct {
	createuc port.CreateBooking
}

func NewBookingHandler(createuc port.CreateBooking) *BookingHandler {
	return &BookingHandler{
		createuc: createuc,
	}
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("tour_id")
	tourID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid agency id", http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)

	payload, err := util.GetPayload(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	userID := payload.UserID
	role := payload.Role

	var result *domain.BookingResponse
	var req domain.BookingRequest
	req.TourID = tourID
	err = decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if role == "user" {

		result, err = h.createuc.Execute(r.Context(), &req, &userID, nil)
	} else if role == "member" {

		result, err = h.createuc.Execute(r.Context(), &req, nil, &userID)
	} else {
		http.Error(w, "Invalid user type", 403)
		return
	}
	if err != nil {
		util.SendData(w, err.Error(), 400)
		return
	}
	util.SendData(w, result, http.StatusCreated)
}
