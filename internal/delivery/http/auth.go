package http

import (
	"net/http"

	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/app_errors"
	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/gorilla/mux"
)

type authController struct {
	authUsecase usecase.IAuthUsecase
}

func NewAuthController(r *mux.Router, authUsecase usecase.IAuthUsecase) {
	controller := &authController{
		authUsecase: authUsecase,
	}
	r.Handle("/register", http.HandlerFunc(controller.Register)).Methods(http.MethodPost)
	r.Handle("/login", http.HandlerFunc(controller.Login)).Methods(http.MethodPost)
	r.Handle("/logout", http.HandlerFunc(controller.Logout)).Methods(http.MethodDelete)
}

// Register регистрирует нового пользователя
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя по данным, переданным в запросе и возвращает token
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body model.RegisterCredentials true "Данные для регистрации пользователя"
// @Success 201 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/register [post]
func (c *authController) Register(w http.ResponseWriter, r *http.Request) {
	var registerValues model.RegisterCredentials

	if err := utils.ParseJSONRequest(r, &registerValues); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid request data", false)
		return
	}

	sessionID, err := c.authUsecase.RegisterUser(r.Context(), registerValues)
	if err != nil {
		code, err := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, code, err.Error(), false)
		return
	}

	utils.SetSessionCookie(w, sessionID)
	utils.SendJSONResponse(w, http.StatusCreated, "Registration successful", true)
}

// Login авторизовывает пользователя
// @Summary Авторизация пользователя
// @Description Авторизовывает, аутентифицирует существующего пользователя и возвращает token
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body model.LoginCredentials true "Данные для авторизации пользователя"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/login [post]
func (c *authController) Login(w http.ResponseWriter, r *http.Request) {
	var loginValues model.LoginCredentials

	if err := utils.ParseJSONRequest(r, &loginValues); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid request data", false)
		return
	}

	sessionID, err := c.authUsecase.LoginUser(r.Context(), loginValues)
	if err != nil {
		code, err := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, code, err.Error(), false)
		return
	}

	utils.SetSessionCookie(w, sessionID)
	utils.SendJSONResponse(w, http.StatusOK, "Login successful", true)
}

// Logout завершает сеанс пользователя
// @Summary Выход пользователя
// @Description Завершает текущую сессию пользователя, удаляя cookie сессии
// @Tags Auth
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/logout [delete]
func (c *authController) Logout(w http.ResponseWriter, r *http.Request) {
	sessionId, err := utils.GetSessionCookie(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "Unauthorized", false)
		return
	}

	if err := c.authUsecase.LogoutUser(r.Context(), sessionId); err != nil {
		code, err := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, code, err.Error(), false)
		return
	}

	utils.DeleteSessionCookie(w)
	utils.SendJSONResponse(w, http.StatusOK, "Logout successful", true)
}
