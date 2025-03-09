package handler

import (
	"net/http"

	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/service"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/utils"
)

// Register регистрирует нового пользователя
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя по данным, переданным в запросе и возвращает token
// @Tags User
// @Accept json
// @Produce json
// @Param user body model.RegisterCredentials true "Данные для регистрации пользователя"
// @Success 201
// @Failure 400
// @Failure 500
// @Router /api/register/ [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var registerValues model.RegisterCredentials

	err := utils.ParseJSONRequest(r, &registerValues)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userResponse, err := service.RegisterUser(registerValues)
	if err != nil {
		http.Error(w, userResponse.Body.(error).Error(), userResponse.StatusCode)
		return
	}

	utils.SetSessionCookie(w, userResponse.Body.(string))

	w.WriteHeader(userResponse.StatusCode)
}

// Login авторизовывает пользователя
// @Summary Авторизация пользователя
// @Description Авторизовывает, аутентифицирует существующего пользователя и возвращает token
// @Tags User
// @Accept json
// @Produce json
// @Param user body model.LoginCredentials true "Данные для авторизации пользователя"
// @Success 201
// @Failure 400
// @Failure 500
// @Router /api/login/ [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var loginValues model.LoginCredentials

	err := utils.ParseJSONRequest(r, &loginValues)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userResponse, err := service.LoginUser(loginValues)
	if err != nil {
		http.Error(w, userResponse.Body.(error).Error(), userResponse.StatusCode)
		return
	}

	utils.SetSessionCookie(w, userResponse.Body.(string))

	w.WriteHeader(userResponse.StatusCode)
}
