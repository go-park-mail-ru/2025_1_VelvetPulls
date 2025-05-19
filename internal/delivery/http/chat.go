package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/app_errors"
	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/proto"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type chatController struct {
	sessionClient authpb.SessionServiceClient
	chatUsecase   usecase.IChatUsecase
}

func NewChatController(r *mux.Router, chatUsecase usecase.IChatUsecase, sessionClient authpb.SessionServiceClient) {
	controller := &chatController{
		chatUsecase:   chatUsecase,
		sessionClient: sessionClient,
	}

	r.Handle("/chats", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.GetChats))).Methods(http.MethodGet)
	r.Handle("/chat", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.CreateChat))).Methods(http.MethodPost)
	r.Handle("/chat/{chat_id}", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.GetChat))).Methods(http.MethodGet)
	r.Handle("/chat/{chat_id}", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.UpdateChat))).Methods(http.MethodPut)
	r.Handle("/chat/{chat_id}", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.DeleteChat))).Methods(http.MethodDelete)
	r.Handle("/chat/{chat_id}/users", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.AddUsersToChat))).Methods(http.MethodPost)
	r.Handle("/chat/{chat_id}/users", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.RemoveUsersFromChat))).Methods(http.MethodDelete)
	r.Handle("/chat/{chat_id}/leave", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.LeaveChat))).Methods(http.MethodPost)
	r.Handle("/chat/{chat_id}/notifications/{send}", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.SendNotifications))).Methods(http.MethodPost)
}

// GetChats возвращает список чатов пользователя
// @Summary Получить список чатов пользователя
// @Description Возвращает список всех чатов, в которых участвует текущий пользователь
// @Tags Chat
// @Produce json
// @Success 200 {array} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chats [get]
func (c *chatController) GetChats(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("GetChats", zap.String("userID", userID.String()))

	chats, err := c.chatUsecase.GetChats(r.Context(), userID)
	if err != nil {
		logger.Error("Failed to get user chats", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, chats, true)
}

// CreateChat создает новый чат
// @Summary Создать новый чат
// @Description Создает новый чат (личный, групповой или канал) с возможностью загрузки аватара
// @Tags Chat
// @Accept multipart/form-data
// @Produce json
// @Param chat_data formData string true "Данные чата в формате JSON"
// @Param avatar formData file false "Аватар чата"
// @Success 201 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat [post]
func (c *chatController) CreateChat(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("CreateChat", zap.String("userID", userID.String()))

	var req model.CreateChatRequest
	if err := utils.ParseJSONRequest(r, &req); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid request body", false)
		return
	}

	chatData := model.CreateChatRequest{
		Type:  req.Type,
		Title: req.Title,
		Users: req.Users,
	}

	chatInfo, err := c.chatUsecase.CreateChat(r.Context(), userID, &chatData)
	if err != nil {
		logger.Error("Failed to create chat", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusCreated, chatInfo, true)
}

// GetChat возвращает информацию о чате
// @Summary Получить информацию о чате
// @Description Возвращает полную информацию о чате по его ID
// @Tags Chat
// @Produce json
// @Param chat_id path string true "ID чата"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 404 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id} [get]
func (c *chatController) GetChat(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	chatID, parseErr := uuid.Parse(mux.Vars(r)["chat_id"])
	if parseErr != nil {
		logger.Error("Invalid chat ID format", zap.Error(parseErr))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat ID", false)
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("GetChat", zap.String("chatID", chatID.String()), zap.String("userID", userID.String()))

	chatInfo, err := c.chatUsecase.GetChatInfo(r.Context(), userID, chatID)
	if err != nil {
		logger.Error("Failed to get chat info", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, chatInfo, true)
}

// @Router /chat/{chat_id}/notifications/{send} [post]
func (c *chatController) SendNotifications(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := utils.GetUserIDFromCtx(r.Context())
	chatID, parseErr := uuid.Parse(mux.Vars(r)["chat_id"])
	if parseErr != nil {
		logger.Error("Invalid chat ID format", zap.Error(parseErr))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat ID", false)
		return
	}

	send, parseErr := strconv.ParseBool(mux.Vars(r)["send"])
	if parseErr != nil {
		logger.Error("Invalid send notification status format", zap.Error(parseErr))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid send bool", false)
		return
	}
	err := c.chatUsecase.SendNotifications(r.Context(), userID, chatID, send)
	if err != nil {
		logger.Error("Failed to set send notifications for chat", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}
	utils.SendJSONResponse(w, r, http.StatusOK, send, true)
}

// UpdateChat обновляет информацию о чате
// @Summary Обновить информацию о чате
// @Description Обновляет информацию о чате (название, аватар) для владельца чата
// @Tags Chat
// @Accept multipart/form-data
// @Produce json
// @Param chat_id path string true "ID чата"
// @Param chat_data formData string true "Данные чата в формате JSON"
// @Param avatar formData file false "Новый аватар чата"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id} [put]
func (c *chatController) UpdateChat(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	chatID, parseErr := uuid.Parse(mux.Vars(r)["chat_id"])
	if parseErr != nil {
		logger.Error("Invalid chat ID format", zap.Error(parseErr))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat ID", false)
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("UpdateChat", zap.String("chatID", chatID.String()), zap.String("userID", userID.String()))

	if err := r.ParseMultipartForm(config.MAX_FILE_SIZE); err != nil {
		logger.Error("Failed to parse multipart form", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Request too large or malformed", false)
		return
	}

	var payload model.UpdateChat
	payload.ID = chatID
	if data := r.FormValue("chat_data"); data != "" {
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			logger.Error("Invalid chat data format", zap.Error(err))
			utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat data format", false)
			return
		}
	}

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

	chatInfo, err := c.chatUsecase.UpdateChat(r.Context(), userID, &payload)
	if err != nil {
		logger.Error("Failed to update chat", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	var updatedAvatar string
	if chatInfo.AvatarPath != nil {
		updatedAvatar = *chatInfo.AvatarPath
	}
	chatInfoResp := model.UpdateChatResp{
		Avatar: updatedAvatar,
		Title:  chatInfo.Title,
	}

	utils.SendJSONResponse(w, r, http.StatusOK, chatInfoResp, true)
}

// DeleteChat удаляет чат
// @Summary Удалить чат
// @Description Удаляет чат (доступно только для владельца чата)
// @Tags Chat
// @Param chat_id path string true "ID чата"
// @Success 204
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id} [delete]
func (c *chatController) DeleteChat(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	chatID, parseErr := uuid.Parse(mux.Vars(r)["chat_id"])
	if parseErr != nil {
		logger.Error("Invalid chat ID format", zap.Error(parseErr))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat ID", false)
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("DeleteChat", zap.String("chatID", chatID.String()), zap.String("userID", userID.String()))

	if err := c.chatUsecase.DeleteChat(r.Context(), userID, chatID); err != nil {
		logger.Error("Failed to delete chat", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddUsersToChat добавляет пользователей в чат
// @Summary Добавить пользователей в чат
// @Description Добавляет одного или нескольких пользователей в чат (доступно для владельца/администратора)
// @Tags Chat
// @Accept json
// @Produce json
// @Param chat_id path string true "ID чата"
// @Param user_ids body []uuid.UUID true "Список ID пользователей для добавления"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id}/users [post]
func (c *chatController) AddUsersToChat(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	chatID, parseErr := uuid.Parse(mux.Vars(r)["chat_id"])
	if parseErr != nil {
		logger.Error("Invalid chat ID format", zap.Error(parseErr))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat ID", false)
		return
	}

	var users model.UsersRequest
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid request body", false)
		return
	}
	usernames := users.Users
	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("AddUsersToChat", zap.String("chatID", chatID.String()), zap.Any("usernames", usernames))

	result, err := c.chatUsecase.AddUsersIntoChat(r.Context(), userID, usernames, chatID)
	if err != nil {
		logger.Error("Failed to add users to chat", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, result, true)
}

// RemoveUsersFromChat удаляет пользователей из чата
// @Summary Удалить пользователей из чата
// @Description Удаляет одного или нескольких пользователей из чата (доступно для владельца/администратора)
// @Tags Chat
// @Accept json
// @Produce json
// @Param chat_id path string true "ID чата"
// @Param user_ids body []uuid.UUID true "Список ID пользователей для удаления"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id}/users [delete]
func (c *chatController) RemoveUsersFromChat(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	chatID, parseErr := uuid.Parse(mux.Vars(r)["chat_id"])
	if parseErr != nil {
		logger.Error("Invalid chat ID format", zap.Error(parseErr))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat ID", false)
		return
	}

	var users model.UsersRequest
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid request body", false)
		return
	}
	usernames := users.Users
	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("RemoveUsersFromChat", zap.String("chatID", chatID.String()), zap.Any("usernames", usernames))

	result, err := c.chatUsecase.DeleteUserFromChat(r.Context(), userID, usernames, chatID)
	if err != nil {
		logger.Error("Failed to remove users from chat", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, result, true)
}

// LeaveChat позволяет текущему пользователю выйти из чата
// @Summary Выйти из чата
// @Description Удаляет текущего пользователя из участников чата (если он не владелец)
// @Tags Chat
// @Produce json
// @Param chat_id path string true "ID чата"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id}/leave [post]
func (c *chatController) LeaveChat(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	chatID, err := uuid.Parse(mux.Vars(r)["chat_id"])
	if err != nil {
		logger.Error("Invalid chat ID format", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat ID", false)
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("LeaveChat", zap.String("chatID", chatID.String()), zap.String("userID", userID.String()))

	if err := c.chatUsecase.LeaveChat(r.Context(), userID, chatID); err != nil {
		logger.Error("Failed to leave chat", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, "left chat successfully", true)
}
