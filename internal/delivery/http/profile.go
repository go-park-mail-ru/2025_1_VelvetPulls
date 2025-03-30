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
	"github.com/google/uuid"
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

	r.Handle("/profile/{user_id}", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.GetProfile))).Methods(http.MethodGet)
	r.Handle("/profile", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.GetSelfProfile))).Methods(http.MethodGet)
	r.Handle("/profile", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.UpdateSelfProfile))).Methods(http.MethodPut)
}

// GetSelfProfile возвращает профиль текущего пользователя
// @Summary Получить профиль текущего пользователя
// @Description Возвращает профиль текущего пользователя, основываясь на ID из контекста сессии
// @Tags User
// @Produce json
// @Success 200 {object} model.UserProfile
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /profile [get]
func (c *userController) GetSelfProfile(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("Fetching self profile")

	profile, err := c.userUsecase.GetUserProfile(r.Context(), userID)
	if err != nil {
		logger.Error("Failed to get self profile", zap.Error(err))
		code, err := apperrors.GetErrAndCodeToSend(err)
		if sendErr := utils.SendJSONResponse(w, code, err, false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusOK, profile, true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
}

// GetProfile возвращает профиль пользователя по ID
// @Summary Получить профиль пользователя по ID
// @Description Возвращает профиль пользователя по предоставленному ID
// @Tags User
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Success 200 {object} model.UserProfile
// @Failure 400 {object} utils.JSONResponse
// @Failure 404 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /profile/{user_id} [get]
func (c *userController) GetProfile(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["user_id"])
	if err != nil {
		logger.Error("Invalid user ID format", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid user ID", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	logger.Info("Fetching user profile", zap.String("userID", userID.String()))
	profile, err := c.userUsecase.GetUserProfile(r.Context(), userID)
	if err != nil {
		logger.Error("Failed to get user profile", zap.Error(err))
		code, err := apperrors.GetErrAndCodeToSend(err)
		if sendErr := utils.SendJSONResponse(w, code, err, false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusOK, profile, true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
}

// UpdateSelfProfile обновляет профиль текущего пользователя
// @Summary Обновить профиль текущего пользователя
// @Description Обновляет профиль текущего пользователя, включая возможность изменить изображение профиля
// @Tags User
// @Accept json
// @Produce json
// @Param profile body model.UserProfile true "Данные профиля"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /profile [put]
func (c *userController) UpdateSelfProfile(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	ctx := r.Context()
	userID := utils.GetUserIDFromCtx(ctx)
	logger.Info("Updating user profile", zap.String("userID", userID.String()))

	if err := r.ParseMultipartForm(config.MAX_FILE_SIZE); err != nil {
		logger.Error("Failed to parse multipart form", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Request too large or malformed", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	var profile model.UpdateUserProfile
	profile.ID = userID

	jsonString := r.FormValue("profile_data")
	if jsonString != "" {
		if err := json.Unmarshal([]byte(jsonString), &profile); err != nil {
			logger.Error("Invalid profile data format", zap.Error(err))
			if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid profile data format", false); sendErr != nil {
				logger.Error("Failed to send error response", zap.Error(sendErr))
			}
			return
		}
	}

	avatar, _, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		logger.Error("Invalid avatar file", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid avatar file", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}
	defer func() {
		if avatar != nil {
			if err := avatar.Close(); err != nil {
				logger.Error("Failed to close avatar file", zap.Error(err))
			}
		}
	}()

	if avatar != nil {
		profile.Avatar = &avatar
	}

	if err := c.userUsecase.UpdateUserProfile(ctx, &profile); err != nil {
		logger.Error("Failed to update user profile", zap.Error(err))
		code, err := apperrors.GetErrAndCodeToSend(err)
		if sendErr := utils.SendJSONResponse(w, code, err, false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusOK, "Profile updated successfully", true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
}
