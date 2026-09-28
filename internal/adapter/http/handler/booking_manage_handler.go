package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/bishal05das/travelbuddy/internal/domain"
	bookingusecase "github.com/bishal05das/travelbuddy/internal/usecase/booking"
	"github.com/bishal05das/travelbuddy/internal/validation"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

// ListAgencyBookings: GET /agency/{agency_id}/bookings?status=&tour_id=&page=&limit=
func (h *BookingHandler) ListAgencyBookings(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorFromRequest(w, r)
	if !ok {
		return
	}
	agencyID, ok := pathUUID(w, r, "agency_id")
	if !ok {
		return
	}
	h.listBookings(w, r, actor, domain.BookingScope{AgencyID: &agencyID})
}

// ListMyBookings: GET /me/bookings?status=&tour_id=&page=&limit=
func (h *BookingHandler) ListMyBookings(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorFromRequest(w, r)
	if !ok {
		return
	}
	h.listBookings(w, r, actor, domain.BookingScope{UserID: &actor.ID})
}

func (h *BookingHandler) listBookings(w http.ResponseWriter, r *http.Request, actor domain.Actor, scope domain.BookingScope) {
	q := r.URL.Query()
	filter := domain.BookingFilter{BookingScope: scope, Status: q.Get("status")}

	var err error
	if filter.Page, err = optionalInt(q.Get("page")); err != nil {
		http.Error(w, "invalid page", http.StatusBadRequest)
		return
	}
	if filter.Limit, err = optionalInt(q.Get("limit")); err != nil {
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}
	if raw := q.Get("tour_id"); raw != "" {
		tourID, err := uuid.Parse(raw)
		if err != nil {
			http.Error(w, "invalid tour_id", http.StatusBadRequest)
			return
		}
		filter.TourID = &tourID
	}

	page, err := h.listUC.Execute(r.Context(), actor, filter)
	if err != nil {
		sendBookingError(w, err)
		return
	}
	util.SendData(w, page, http.StatusOK)
}

// GetAgencyBooking: GET /agency/{agency_id}/bookings/{booking_id}
func (h *BookingHandler) GetAgencyBooking(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorFromRequest(w, r)
	if !ok {
		return
	}
	agencyID, ok := pathUUID(w, r, "agency_id")
	if !ok {
		return
	}
	h.getBooking(w, r, actor, domain.BookingScope{AgencyID: &agencyID})
}

// GetMyBooking: GET /me/bookings/{booking_id}
func (h *BookingHandler) GetMyBooking(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorFromRequest(w, r)
	if !ok {
		return
	}
	h.getBooking(w, r, actor, domain.BookingScope{UserID: &actor.ID})
}

func (h *BookingHandler) getBooking(w http.ResponseWriter, r *http.Request, actor domain.Actor, scope domain.BookingScope) {
	bookingID, ok := pathUUID(w, r, "booking_id")
	if !ok {
		return
	}
	booking, err := h.getUC.Execute(r.Context(), actor, scope, bookingID)
	if err != nil {
		sendBookingError(w, err)
		return
	}
	util.SendData(w, booking, http.StatusOK)
}

// UpdateAgencyBookingStatus: PATCH /agency/{agency_id}/bookings/{booking_id}/status
// with {"status": "confirmed" | "cancelled" | "completed"}.
func (h *BookingHandler) UpdateAgencyBookingStatus(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorFromRequest(w, r)
	if !ok {
		return
	}
	agencyID, ok := pathUUID(w, r, "agency_id")
	if !ok {
		return
	}
	bookingID, ok := pathUUID(w, r, "booking_id")
	if !ok {
		return
	}
	var req domain.UpdateBookingStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := validation.Validate.Struct(req); err != nil {
		http.Error(w, "status must be one of: confirmed, cancelled, completed", http.StatusBadRequest)
		return
	}

	booking, err := h.updateStatusUC.Execute(r.Context(), actor, agencyID, bookingID, req.Status)
	if err != nil {
		sendBookingError(w, err)
		return
	}
	util.SendData(w, booking, http.StatusOK)
}

// CancelMyBooking: POST /me/bookings/{booking_id}/cancel
func (h *BookingHandler) CancelMyBooking(w http.ResponseWriter, r *http.Request) {
	actor, ok := actorFromRequest(w, r)
	if !ok {
		return
	}
	bookingID, ok := pathUUID(w, r, "booking_id")
	if !ok {
		return
	}
	booking, err := h.cancelUC.Execute(r.Context(), actor, bookingID)
	if err != nil {
		sendBookingError(w, err)
		return
	}
	util.SendData(w, booking, http.StatusOK)
}

func sendBookingError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrBookingNotFound):
		util.SendData(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, bookingusecase.ErrInvalidStatusFilter):
		util.SendData(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, bookingusecase.ErrForbiddenScope):
		util.SendData(w, err.Error(), http.StatusForbidden)
	case errors.Is(err, domain.ErrInvalidBookingTransition), errors.Is(err, bookingusecase.ErrTourAlreadyStarted):
		util.SendData(w, err.Error(), http.StatusConflict)
	default:
		log.Println("booking request failed:", err)
		util.SendData(w, "internal server error", http.StatusInternalServerError)
	}
}

func pathUUID(w http.ResponseWriter, r *http.Request, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(r.PathValue(name))
	if err != nil {
		http.Error(w, "invalid "+name, http.StatusBadRequest)
		return uuid.Nil, false
	}
	return id, true
}

func optionalInt(v string) (int, error) {
	if v == "" {
		return 0, nil
	}
	return strconv.Atoi(v)
}
