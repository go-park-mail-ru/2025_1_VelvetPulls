package http

import (
	"encoding/json"
	"net/http"

	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
	userID, err := uuid.Parse(utils.GetUserIDFromCtx(r.Context()))
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid userID format", false)
		return
	}

	profile, err := c.userUsecase.GetUserProfile(r.Context(), userID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), false)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, profile, true)
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
	vars := mux.Vars(r)
	userID, err := uuid.Parse(vars["user_id"])
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid userID format", false)
		return
	}

	profile, err := c.userUsecase.GetUserProfile(r.Context(), userID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), false)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, profile, true)
}

// UpdateSelfProfile обновляет профиль текущего пользователя
// @Summary Обновить профиль текущего пользователя
// @Description Обновляет профиль текущего пользователя, включая возможность изменить изображение профиля
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Param profile body model.UserProfile true "Данные профиля"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /profile [put]
func (c *userController) UpdateSelfProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := uuid.Parse(utils.GetUserIDFromCtx(ctx))
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid userID format", false)
		return
	}
	var profile model.UpdateUserProfile

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "Unable to parse form", false)
		return
	}
	profile.ID = userID

	jsonString := r.FormValue("profile_data")
	if err := json.Unmarshal([]byte(jsonString), &profile); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid JSON format in profile_data", false)
		return
	}

	avatar, _, err := r.FormFile("avatar")
	if err != nil && err != http.ErrMissingFile {
		utils.SendJSONResponse(w, http.StatusBadRequest, "Failed to get avatar file", false)
		return
	}
	defer func() {
		if avatar != nil {
			avatar.Close()
		}
	}()

	if avatar != nil {
		profile.Avatar = &avatar
	}

	if err := c.userUsecase.UpdateUserProfile(ctx, &profile); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), false)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "Profile updated successfully", true)
}
