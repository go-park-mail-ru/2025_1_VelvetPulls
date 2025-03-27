package http

import (
	"net/http"
	"strconv"

	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
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
	r.Handle("/profile/{user_id}", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.UpdateProfile))).Methods(http.MethodPut)
	r.Handle("/avatar/{user_id}", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.GetAvatar))).Methods(http.MethodGet)
	r.Handle("/avatar/{user_id}", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.UpdateAvatar))).Methods(http.MethodPut)
}

// GetProfile возвращает профиль пользователя
// @Summary Получить профиль пользователя
// @Description Возвращает профиль пользователя по ID
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
	userID := vars["user_id"]

	profile, err := c.userUsecase.GetUserProfile(r.Context(), userID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), false)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, profile, true)
}

// UpdateProfile обновляет профиль пользователя
// @Summary Обновить профиль пользователя
// @Description Обновляет данные профиля пользователя
// @Tags User
// @Accept json
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Param profile body model.UserProfile true "Данные профиля"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /profile/{user_id} [put]
func (c *userController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var profile model.UserProfile
	if err := utils.ParseJSONRequest(r, &profile); err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid request body", false)
		return
	}

	if err := c.userUsecase.UpdateUserProfile(ctx, &profile); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), false)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "Profile updated successfully", true)
}

// GetAvatar возвращает аватар пользователя
// @Summary Получить аватар пользователя
// @Description Возвращает аватар пользователя по ID
// @Tags User
// @Produce octet-stream
// @Param user_id path string true "ID пользователя"
// @Success 200 {file} binary
// @Failure 404 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /avatar/{user_id} [get]
func (c *userController) GetAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	userID := vars["user_id"]

	avatarBytes, err := c.userUsecase.GetUserAvatar(ctx, userID)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, "Failed to retrieve avatar", false)
		return
	}

	if avatarBytes == nil {
		utils.SendJSONResponse(w, http.StatusNotFound, "Avatar not found", false)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(avatarBytes)))
	if _, err := w.Write(avatarBytes); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, "Failed to send avatar", false)
	}
}

// UpdateAvatar обновляет аватар пользователя
// @Summary Обновить аватар пользователя
// @Description Загружает новый аватар для пользователя
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param user_id path string true "ID пользователя"
// @Param avatar formData file true "Файл аватара"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /avatar/{user_id} [put]
func (c *userController) UpdateAvatar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	file, header, err := r.FormFile("avatar")
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "Failed to get avatar from form", false)
		return
	}
	defer file.Close()

	if err := c.userUsecase.UploadAvatar(ctx, file, header); err != nil {
		utils.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), false)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, "Avatar updated successfully", true)
}
