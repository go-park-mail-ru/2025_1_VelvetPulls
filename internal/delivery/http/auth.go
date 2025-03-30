package http

import (
	"net/http"

	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/app_errors"
	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	usecase "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
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
	logger := utils.GetLoggerFromCtx(r.Context())

	var registerValues model.RegisterCredentials
	if err := utils.ParseJSONRequest(r, &registerValues); err != nil {
		logger.Warn("Invalid request data", zap.Error(err))
		if err := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid request data", false); err != nil {
			logger.Error("Failed to send JSON response", zap.Error(err))
		}
		return
	}

	sessionID, err := c.authUsecase.RegisterUser(r.Context(), registerValues)
	if err != nil {
		code, errMsg := apperrors.GetErrAndCodeToSend(err)
		logger.Error("Registration failed", zap.String("error", errMsg.Error()))
		if err := utils.SendJSONResponse(w, code, errMsg, false); err != nil {
			logger.Error("Failed to send JSON response", zap.Error(err))
		}
		return
	}

	utils.SetSessionCookie(w, sessionID)
	logger.Info("User registered successfully")
	if err := utils.SendJSONResponse(w, http.StatusCreated, "Registration successful", true); err != nil {
		logger.Error("Failed to send JSON response", zap.Error(err))
	}
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
	logger := utils.GetLoggerFromCtx(r.Context())

	var loginValues model.LoginCredentials
	if err := utils.ParseJSONRequest(r, &loginValues); err != nil {
		logger.Warn("Invalid request data", zap.Error(err))
		if err := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid request data", false); err != nil {
			logger.Error("Failed to send JSON response", zap.Error(err))
		}
		return
	}

	sessionID, err := c.authUsecase.LoginUser(r.Context(), loginValues)
	if err != nil {
		code, errMsg := apperrors.GetErrAndCodeToSend(err)
		logger.Error("Login failed", zap.String("error", errMsg.Error()))
		if err := utils.SendJSONResponse(w, code, errMsg, false); err != nil {
			logger.Error("Failed to send JSON response", zap.Error(err))
		}
		return
	}

	utils.SetSessionCookie(w, sessionID)
	logger.Info("User logged in successfully")
	if err := utils.SendJSONResponse(w, http.StatusOK, "Login successful", true); err != nil {
		logger.Error("Failed to send JSON response", zap.Error(err))
	}
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
	logger := utils.GetLoggerFromCtx(r.Context())

	sessionId, err := utils.GetSessionCookie(r)
	if err != nil {
		logger.Warn("Unauthorized logout attempt")
		if err := utils.SendJSONResponse(w, http.StatusBadRequest, "Unauthorized", false); err != nil {
			logger.Error("Failed to send JSON response", zap.Error(err))
		}
		return
	}

	if err := c.authUsecase.LogoutUser(r.Context(), sessionId); err != nil {
		code, errMsg := apperrors.GetErrAndCodeToSend(err)
		logger.Error("Logout failed", zap.String("error", errMsg.Error()))
		if err := utils.SendJSONResponse(w, code, errMsg, false); err != nil {
			logger.Error("Failed to send JSON response", zap.Error(err))
		}
		return
	}

	utils.DeleteSessionCookie(w)
	logger.Info("User logged out successfully")
	if err := utils.SendJSONResponse(w, http.StatusOK, "Logout successful", true); err != nil {
		logger.Error("Failed to send JSON response", zap.Error(err))
	}
}
