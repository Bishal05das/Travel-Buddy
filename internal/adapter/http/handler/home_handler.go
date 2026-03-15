package handler

import (
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
)

type HomeHandler struct {
	uc port.Home
}

func NewHomeHandler(uc port.Home) *HomeHandler {
	return &HomeHandler{uc: uc}
}

func (h *HomeHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		util.SendData(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp, err := h.uc.GetHome(r.Context())
	if err != nil {
		util.SendData(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	util.SendData(w, map[string]any{
		"success": true,
		"data":    resp,
	},http.StatusOK)
}
