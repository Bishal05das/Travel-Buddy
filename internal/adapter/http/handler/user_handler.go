package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
	util "github.com/bishal05das/travelbuddy/utils"
)

type UserHandler struct {
	createuc port.CreateUser
	loginuc port.LoginUser
}

func NewUserHandler(createuc port.CreateUser,loginuc port.LoginUser) *UserHandler {
	return &UserHandler{
		createuc: createuc,
		loginuc: loginuc,
	}
}

func(h *UserHandler) CreateUser(w http.ResponseWriter,r *http.Request) {
	var newUser domain.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&newUser)
	if err != nil {
		http.Error(w,"Invalid Request Data",http.StatusBadRequest)
		return
	}
	err = h.createuc.Execute(r.Context(),&newUser)
	if err != nil {
		http.Error(w,err.Error(), http.StatusInternalServerError)
		return
	}
	util.SendData(w,newUser,http.StatusCreated)

}

func(h *UserHandler) UserLogin(w http.ResponseWriter,r *http.Request) {
	var reqLogin domain.ReqLogin

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqLogin)
	if err != nil {
		http.Error(w,err.Error(),http.StatusBadRequest)
		return
	}
	token,err := h.loginuc.Execute(r.Context(),&reqLogin)
	if err != nil {
		util.SendData(w,err.Error(),http.StatusInternalServerError)
		return
	}
	util.SendData(w,token,http.StatusCreated)
}