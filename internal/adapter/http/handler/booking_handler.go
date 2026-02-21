package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
)

type BookingHandler struct {
	createuc port.CreateBooking
}

func NewBookingHandler(createuc port.CreateBooking) *BookingHandler {
	return &BookingHandler{
		createuc: createuc,
	}
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter,r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	payload,err := util.GetPayload(r)
	if err != nil {
		http.Error(w,err.Error(),http.StatusUnauthorized)
	}
	userID := payload.UserID
	role := payload.Role

	var result *domain.BookingResponse
	if role == "user" {
		var req domain.BookingRequest
		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w,err.Error(),http.StatusBadRequest)
		}
		result,err = h.createuc.Execute(r.Context(),&req,&userID,nil)
	}else if role == "admin" || role == "sub" {
		var req domain.BookingRequest
		err := decoder.Decode(&req)
		if err != nil {
			http.Error(w,err.Error(),http.StatusBadRequest)
		}
		result,err = h.createuc.Execute(r.Context(),&req,nil,&userID)
	}else {
		http.Error(w,"Invalid user type",403)
	}
	if err != nil {
		util.SendData(w,err.Error(),400)
		return
	}
	util.SendData(w,result,http.StatusCreated)
}