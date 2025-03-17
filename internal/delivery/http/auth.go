package http

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/utils"
	"github.com/gorilla/mux"
)

type authController struct {
	authUsecase usecase.IAuthUsecase
}

func NewAuthController(r *mux.Router, authUsecase usecase.IAuthUsecase) {
	controller := &authController{
		authUsecase: authUsecase,
	}

	r.HandleFunc("/register/", controller.Register).Methods(http.MethodPost)
	r.HandleFunc("/login/", controller.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout/", controller.Logout).Methods(http.MethodDelete)
}

// Register регистрирует нового пользователя
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя по данным, переданным в запросе и возвращает token
// @Tags User
// @Accept json
// @Produce json
// @Param user body model.RegisterCredentials true "Данные для регистрации пользователя"
// @Success 201 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/register/ [post]
func (c *authController) Register(w http.ResponseWriter, r *http.Request) {
	var registerValues model.RegisterCredentials

	// Парсим JSON из запроса
	err := utils.ParseJSONRequest(r, &registerValues)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), false)
		return
	}

	// Регистрируем пользователя
	sessionID, err := c.authUsecase.RegisterUser(registerValues)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrPasswordsDoNotMatch),
			errors.Is(err, apperrors.ErrInvalidPassword),
			errors.Is(err, apperrors.ErrInvalidPhoneFormat),
			errors.Is(err, apperrors.ErrInvalidUsername):
			utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), false)

		case errors.Is(err, apperrors.ErrUsernameTaken),
			errors.Is(err, apperrors.ErrEmailTaken),
			errors.Is(err, apperrors.ErrPhoneTaken):
			utils.SendJSONResponse(w, http.StatusConflict, err.Error(), false)

		default:
			utils.SendJSONResponse(w, http.StatusInternalServerError, apperrors.ErrInternalServer, false)
		}
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
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/login/ [post]
func (c *authController) Login(w http.ResponseWriter, r *http.Request) {
	var loginValues model.LoginCredentials

	// Парсим JSON из запроса
	err := utils.ParseJSONRequest(r, &loginValues)
	fmt.Print(loginValues)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), false)
		return
	}

	// Авторизация пользователя
	sessionID, err := c.authUsecase.LoginUser(loginValues)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUserNotFound),
			errors.Is(err, apperrors.ErrInvalidCredentials):
			utils.SendJSONResponse(w, http.StatusUnauthorized, err.Error(), false)

		default:
			utils.SendJSONResponse(w, http.StatusInternalServerError, apperrors.ErrInternalServer, false)
		}
		return
	}

	// Устанавливаем cookie сессии
	utils.SetSessionCookie(w, sessionID)

	// Отправляем успешный ответ
	utils.SendJSONResponse(w, http.StatusOK, "Login successful", true)
}

// Logout завершает сеанс пользователя
// @Summary Выход пользователя
// @Description Завершает текущую сессию пользователя, удаляя cookie сессии
// @Tags User
// @Success 200 {object} utils.JSONResponse
// @Router /api/logout/ [delete]
func (c *authController) Logout(w http.ResponseWriter, r *http.Request) {
	utils.DeleteSessionCookie(w)
	utils.SendJSONResponse(w, http.StatusOK, "Logout successful", true)
}
