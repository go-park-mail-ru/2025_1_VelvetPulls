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
// @Success 201 {string} string
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/register/ [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var registerValues model.RegisterCredentials

	// Парсим JSON из запроса
	err := utils.ParseJSONRequest(r, &registerValues)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), false)
		return
	}

	// Регистрируем пользователя
	sessionID, err := service.RegisterUser(registerValues)
	if err != nil {
		// Отправляем ошибку в формате JSON
		utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), false)
		return
	}

	// Устанавливаем cookie сессии
	utils.SetSessionCookie(w, sessionID)

	// Отправляем успешный ответ
	utils.SendJSONResponse(w, http.StatusCreated, "Registration successful", true)
}

// Login авторизовывает пользователя
// @Summary Авторизация пользователя
// @Description Авторизовывает, аутентифицирует существующего пользователя и возвращает token
// @Tags User
// @Accept json
// @Produce json
// @Param user body model.LoginCredentials true "Данные для авторизации пользователя"
// @Success 200 {string} string
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/login/ [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var loginValues model.LoginCredentials

	// Парсим JSON из запроса
	err := utils.ParseJSONRequest(r, &loginValues)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), false)
		return
	}

	// Авторизация пользователя
	sessionID, err := service.LoginUser(loginValues)
	if err != nil {
		// Отправляем ошибку в формате JSON
		utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), false)
		return
	}

	// Устанавливаем cookie сессии
	utils.SetSessionCookie(w, sessionID)

	// Отправляем успешный ответ
	utils.SendJSONResponse(w, http.StatusOK, "Login successful", true)
}
