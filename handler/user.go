package handler

import (
	"net/http"

	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	service "github.com/go-park-mail-ru/2025_1_VelvetPulls/service"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/utils"
	"github.com/gorilla/mux"
)

type UserService interface {
	RegisterUser(user model.User) (service.UserResponse, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserHandler(r *mux.Router, userService UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := utils.ParseJSONRequest(r, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userResponse, err := h.userService.RegisterUser(user)
	if err != nil {
		http.Error(w, userResponse.Body.(error).Error(), userResponse.StatusCode)
		return
	}

	err = utils.SendJSONResponse(w, userResponse.StatusCode, userResponse.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
