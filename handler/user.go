package handler

import (
	"net/http"

	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/service"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/utils"
)

// Register регистрирует нового пользователя
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя по данным, переданным в запросе
// @Tags User
// @Accept json
// @Produce json
// @Param user body model.User true "Данные для регистрации пользователя"
// @Success 201 {object} model.User
// @Failure 400
// @Failure 500
// @Router /api/register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := utils.ParseJSONRequest(r, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userResponse, err := service.RegisterUser(user)
	if err != nil {
		http.Error(w, userResponse.Body.(error).Error(), userResponse.StatusCode)
		return
	}

	err = utils.SendJSONResponse(w, userResponse.StatusCode, userResponse.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Login авторизовывает пользователя
// @Summary Авторизация пользователя
// @Description Авторизовывает и аутентифицирует существующего пользователя, создавая новую сессию
// @Tags User Session
// @Accept json
// @Produce json
// @Param user body model.AuthCredentials true "Данные для авторизации пользователя"
// @Success 201 {object} model.Session
// @Failure 400
// @Failure 500
// @Router /api/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var loginValues model.AuthCredentials

	err := utils.ParseJSONRequest(r, &loginValues)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var newSession model.Session
	userResponse, err := service.AuthenticateUser(loginValues, newSession)
	if err != nil {
		http.Error(w, userResponse.Body.(error).Error(), userResponse.StatusCode)
		return
	}
	sessionId := userResponse.Body.(map[string]interface{})["username"].(string)

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	err = utils.SendJSONResponse(w, userResponse.StatusCode, userResponse.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
