package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/app_errors"
	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type userController struct {
	sessionUsecase usecase.ISessionUsecase
	userUsecase    usecase.IUserUsecase
}

func NewUserController(r *mux.Router, userUsecase usecase.IUserUsecase, sessionUsecase usecase.ISessionUsecase) {
	controller := &userController{
		userUsecase:    userUsecase,
		sessionUsecase: sessionUsecase,
	}

	r.Handle("/profile/{username}", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.GetProfile))).Methods(http.MethodGet)
	r.Handle("/profile", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.GetSelfProfile))).Methods(http.MethodGet)
	r.Handle("/profile", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.UpdateSelfProfile))).Methods(http.MethodPut)
}

// GetSelfProfile возвращает профиль текущего пользователя
// @Summary Получить профиль текущего пользователя
// @Description Возвращает профиль текущего пользователя, основываясь на ID из контекста сессии
// @Tags User
// @Produce json
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /profile [get]
func (c *userController) GetSelfProfile(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := utils.GetUserIDFromCtx(r.Context())

	logger.Info("GetSelfProfile")

	profile, err := c.userUsecase.GetUserProfileByID(r.Context(), userID)
	if err != nil {
		logger.Error("Failed to get self profile", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, profile, true)
}

// GetProfile возвращает профиль пользователя по ID
// @Summary Получить профиль пользователя по ID
// @Description Возвращает профиль пользователя по предоставленному ID
// @Tags User
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 404 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /profile/{username} [get]
func (c *userController) GetProfile(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	username := mux.Vars(r)["username"]

	logger.Info("GetProfile", zap.String("username", username))

	profile, err := c.userUsecase.GetUserProfileByUsername(r.Context(), username)
	if err != nil {
		logger.Error("Failed to get user profile", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, profile, true)
}

// UpdateSelfProfile обновляет профиль текущего пользователя
// @Summary Обновить профиль текущего пользователя
// @Description Обновляет профиль текущего пользователя, включая возможность изменить изображение профиля
// @Tags User
// @Accept json
// @Produce json
// @Param profile body model.UpdateUserProfile true "Данные профиля"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /profile [put]
func (c *userController) UpdateSelfProfile(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := utils.GetUserIDFromCtx(r.Context())

	logger.Info("UpdateSelfProfile", zap.String("userID", userID.String()))

	// 1) Multipart parsing
	if err := r.ParseMultipartForm(config.MAX_FILE_SIZE); err != nil {
		logger.Error("Failed to parse multipart form", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Request too large or malformed", false)
		return
	}

	// 2) JSON part
	var payload model.UpdateUserProfile
	payload.ID = userID
	if data := r.FormValue("profile_data"); data != "" {
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			logger.Error("Invalid profile data format", zap.Error(err))
			utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid profile data format", false)
			return
		}
	}

	// 3) Optional avatar
	avatar, _, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		logger.Error("Invalid avatar file", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid avatar file", false)
		return
	}
	if avatar != nil {
		defer avatar.Close()
		payload.Avatar = &avatar
	}

	// 4) Business update
	if err := c.userUsecase.UpdateUserProfile(r.Context(), &payload); err != nil {
		logger.Error("Failed to update user profile", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, "Profile updated successfully", true)
}
