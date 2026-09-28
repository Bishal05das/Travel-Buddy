package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/bishal05das/travelbuddy/internal/validation"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

type BookingRequest[T any] interface {
	*T
	SetTourID(id uuid.UUID)
	ToCommand(actorID *uuid.UUID) *domain.BookingCommand
}

type BookingHandler struct {
	createuc port.CreateBooking
}

func NewBookingHandler(createuc port.CreateBooking) *BookingHandler {
	return &BookingHandler{
		createuc: createuc,
	}
}

func (h *BookingHandler) CreateBookingByUser(w http.ResponseWriter, r *http.Request) {
	handleCreateBooking[domain.BookingRequestByUser](
		h, w, r, "user",
	)
}

func (h *BookingHandler) CreateBookingByAdmin(w http.ResponseWriter, r *http.Request) {
	// Guest bookings are made by agency members; tokens never carry an
	// "admin" role, which made this endpoint unreachable.
	handleCreateBooking[domain.BookingRequestByAdmin](
		h, w, r, "member",
	)
}

func handleCreateBooking[T any, PT BookingRequest[T]](h *BookingHandler, w http.ResponseWriter, r *http.Request, requiredRole string) {
	idStr := r.PathValue("tour_id")
	tourID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid tour id", http.StatusBadRequest)
		return
	}
	payload, err := util.GetPayload(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	actorID := payload.UserID

	if payload.Role != requiredRole {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req T
	pReq := PT(&req)

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(pReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pReq.SetTourID(tourID)

	if err := validation.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cmd := pReq.ToCommand(&actorID)
	cmd.AgencyID = payload.AgencyID
	result, err := h.createuc.Execute(r.Context(), cmd)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, result, http.StatusCreated)

}
