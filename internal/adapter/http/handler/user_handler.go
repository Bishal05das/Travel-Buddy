package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	"github.com/bishal05das/travelbuddy/internal/validation"
	util "github.com/bishal05das/travelbuddy/utils"
	"github.com/google/uuid"
)

type UserHandler struct {
	createUC port.CreateUser
	loginUC  port.LoginUser
	deleteUC port.DeleteUser
	updateUC port.UpdateUser
}

func NewUserHandler(createUC port.CreateUser, loginUC port.LoginUser, deleteUC port.DeleteUser, updateUC port.UpdateUser) *UserHandler {
	return &UserHandler{
		createUC: createUC,
		loginUC:  loginUC,
		deleteUC: deleteUC,
		updateUC: updateUC,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req domain.CreateUserReq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request Data", http.StatusBadRequest)
		return
	}
	if err := validation.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newUser := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
	}
	err = h.createUC.Execute(r.Context(), &newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	util.SendData(w, "Successfully Created User", http.StatusCreated)

}

func (h *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var reqLogin domain.ReqLogin

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = validation.Validate.Struct(reqLogin); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := h.loginUC.Execute(r.Context(), &reqLogin)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusInternalServerError)
		return
	}
	util.SendData(w, token, http.StatusCreated)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("user_id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid agency id", http.StatusBadRequest)
		return
	}
	err = h.deleteUC.Execute(r.Context(), userID)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, "Successfully Deleted User", http.StatusCreated)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("user_id")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid agency id", http.StatusBadRequest)
		return
	}
	var req domain.UpdateUserReq
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		http.Error(w, "Error in decoding request body", http.StatusBadRequest)
		return
	}
	if err = validation.Validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := &domain.User{
		UserID:    userID,
		Name:      req.Name,
		Email:     req.Email,
		Password:  req.Password,
		Phone:     req.Phone,
		UpdatedAt: time.Now(),
	}

	err = h.updateUC.Execute(r.Context(), user)
	if err != nil {
		util.SendData(w, err.Error(), http.StatusBadRequest)
		return
	}
	util.SendData(w, "Succesfully Updated User", http.StatusOK)
}
