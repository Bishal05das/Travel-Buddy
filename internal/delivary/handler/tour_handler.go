package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/domain"
	tourusecase "github.com/bishal05das/travelbuddy/internal/usecase/tour"
	util "github.com/bishal05das/travelbuddy/utils"
)

type TourHandler struct {
	createUC *tourusecase.CreateTourUseCase
}

func NewTourHandler(createUC *tourusecase.CreateTourUseCase) *TourHandler {
	return &TourHandler{
		createUC: createUC,
	}
}

func (h *TourHandler) Create(w http.ResponseWriter,r *http.Request) {
	var tour domain.Tour
	err := json.NewDecoder(r.Body).Decode(&tour)
	if err != nil {
		fmt.Println("Error: ",err)
		return
	}
	h.createUC.Execute(&tour)
	util.SendDate(w,tour,http.StatusCreated)
}