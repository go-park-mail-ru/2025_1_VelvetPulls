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

type chatController struct {
	sessionUsecase usecase.ISessionUsecase
	chatUsecase    usecase.IChatUsecase
}

func NewChatController(r *mux.Router, chatUsecase usecase.IChatUsecase, sessionUsecase usecase.ISessionUsecase) {
	controller := &chatController{
		chatUsecase:    chatUsecase,
		sessionUsecase: sessionUsecase,
	}

	r.Handle("/chats", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.GetChats))).Methods(http.MethodGet)
	r.Handle("/chat", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.CreateChat))).Methods(http.MethodPost)
	r.Handle("/chat/{chat_id}", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.GetChat))).Methods(http.MethodGet)
	r.Handle("/chat/{chat_id}", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.UpdateChat))).Methods(http.MethodPut)
	r.Handle("/chat/{chat_id}", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.DeleteChat))).Methods(http.MethodDelete)
	r.Handle("/chat/{chat_id}/users", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.AddUsersToChat))).Methods(http.MethodPost)
	r.Handle("/chat/{chat_id}/users", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.RemoveUsersFromChat))).Methods(http.MethodDelete)
	r.Handle("/chat/{chat_id}/leave", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.LeaveChat))).Methods(http.MethodPost)
}

// GetChats возвращает список чатов пользователя
// @Summary Получить список чатов пользователя
// @Description Возвращает список всех чатов, в которых участвует текущий пользователь
// @Tags Chat
// @Produce json
// @Success 200 {array} model.Chat
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chats [get]
func (c *chatController) GetChats(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("Fetching user chats")

	chats, err := c.chatUsecase.GetChats(r.Context(), userID)
	if err != nil {
		logger.Error("Failed to get user chats", zap.Error(err))
		code, err := apperrors.GetErrAndCodeToSend(err)
		if sendErr := utils.SendJSONResponse(w, code, err, false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusOK, chats, true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
}

// CreateChat создает новый чат
// @Summary Создать новый чат
// @Description Создает новый чат (личный, групповой или канал) с возможностью загрузки аватара
// @Tags Chat
// @Accept multipart/form-data
// @Produce json
// @Param chat_data formData string true "Данные чата в формате JSON"
// @Param avatar formData file false "Аватар чата"
// @Success 201 {object} model.ChatInfo
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat [post]
func (c *chatController) CreateChat(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	ctx := r.Context()
	userID := utils.GetUserIDFromCtx(ctx)
	logger.Info("Creating new chat")

	if err := r.ParseMultipartForm(config.MAX_FILE_SIZE); err != nil {
		logger.Error("Failed to parse multipart form", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Request too large or malformed", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	var chatData model.CreateChat
	jsonString := r.FormValue("chat_data")
	if jsonString != "" {
		if err := json.Unmarshal([]byte(jsonString), &chatData); err != nil {
			logger.Error("Invalid chat data format", zap.Error(err))
			if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid chat data format", false); sendErr != nil {
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
		chatData.Avatar = &avatar
	}

	chatInfo, err := c.chatUsecase.CreateChat(ctx, userID, &chatData)
	if err != nil {
		logger.Error("Failed to create chat", zap.Error(err))
		code, err := apperrors.GetErrAndCodeToSend(err)
		if sendErr := utils.SendJSONResponse(w, code, err, false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusCreated, chatInfo, true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
}

// GetChat возвращает информацию о чате
// @Summary Получить информацию о чате
// @Description Возвращает полную информацию о чате по его ID
// @Tags Chat
// @Produce json
// @Param chat_id path string true "ID чата"
// @Success 200 {object} model.ChatInfo
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 404 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id} [get]
func (c *chatController) GetChat(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	vars := mux.Vars(r)
	chatID, err := uuid.Parse(vars["chat_id"])
	if err != nil {
		logger.Error("Invalid chat ID format", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid chat ID", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("Fetching chat info", zap.String("chatID", chatID.String()))

	chatInfo, err := c.chatUsecase.GetChatInfo(r.Context(), userID, chatID)
	if err != nil {
		logger.Error("Failed to get chat info", zap.Error(err))
		code, err := apperrors.GetErrAndCodeToSend(err)
		if sendErr := utils.SendJSONResponse(w, code, err, false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusOK, chatInfo, true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
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
// @Success 200 {object} model.ChatInfo
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id} [put]
func (c *chatController) UpdateChat(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	ctx := r.Context()
	vars := mux.Vars(r)
	chatID, err := uuid.Parse(vars["chat_id"])
	if err != nil {
		logger.Error("Invalid chat ID format", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid chat ID", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	userID := utils.GetUserIDFromCtx(ctx)
	logger.Info("Updating chat", zap.String("chatID", chatID.String()))

	if err := r.ParseMultipartForm(config.MAX_FILE_SIZE); err != nil {
		logger.Error("Failed to parse multipart form", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Request too large or malformed", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	var chatData model.UpdateChat
	chatData.ID = chatID

	jsonString := r.FormValue("chat_data")
	if jsonString != "" {
		if err := json.Unmarshal([]byte(jsonString), &chatData); err != nil {
			logger.Error("Invalid chat data format", zap.Error(err))
			if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid chat data format", false); sendErr != nil {
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
		chatData.Avatar = &avatar
	}

	chatInfo, err := c.chatUsecase.UpdateChat(ctx, userID, &chatData)
	if err != nil {
		logger.Error("Failed to update chat", zap.Error(err))
		code, err := apperrors.GetErrAndCodeToSend(err)
		if sendErr := utils.SendJSONResponse(w, code, err, false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusOK, chatInfo, true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
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
	vars := mux.Vars(r)
	chatID, err := uuid.Parse(vars["chat_id"])
	if err != nil {
		logger.Error("Invalid chat ID format", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid chat ID", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("Deleting chat", zap.String("chatID", chatID.String()))

	if err := c.chatUsecase.DeleteChat(r.Context(), userID, chatID); err != nil {
		logger.Error("Failed to delete chat", zap.Error(err))
		code, err := apperrors.GetErrAndCodeToSend(err)
		if sendErr := utils.SendJSONResponse(w, code, err, false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
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
// @Success 200 {object} model.AddedUsersIntoChat
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id}/users [post]
func (c *chatController) AddUsersToChat(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	vars := mux.Vars(r)
	chatID, err := uuid.Parse(vars["chat_id"])
	if err != nil {
		logger.Error("Invalid chat ID format", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid chat ID", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	var userIDs []uuid.UUID
	if err := json.NewDecoder(r.Body).Decode(&userIDs); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid request body", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("Adding users to chat",
		zap.String("chatID", chatID.String()),
		zap.Any("userIDs", userIDs))

	result, err := c.chatUsecase.AddUsersIntoChat(r.Context(), userID, userIDs, chatID)
	if err != nil {
		logger.Error("Failed to add users to chat", zap.Error(err))
		code, err := apperrors.GetErrAndCodeToSend(err)
		if sendErr := utils.SendJSONResponse(w, code, err, false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusOK, result, true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
}

// RemoveUsersFromChat удаляет пользователей из чата
// @Summary Удалить пользователей из чата
// @Description Удаляет одного или нескольких пользователей из чата (доступно для владельца/администратора)
// @Tags Chat
// @Accept json
// @Produce json
// @Param chat_id path string true "ID чата"
// @Param user_ids body []uuid.UUID true "Список ID пользователей для удаления"
// @Success 200 {object} model.DeletedUsersFromChat
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id}/users [delete]
func (c *chatController) RemoveUsersFromChat(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	vars := mux.Vars(r)
	chatID, err := uuid.Parse(vars["chat_id"])
	if err != nil {
		logger.Error("Invalid chat ID format", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid chat ID", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	var userIDs []uuid.UUID
	if err := json.NewDecoder(r.Body).Decode(&userIDs); err != nil {
		logger.Error("Invalid request body", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid request body", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("Removing users from chat",
		zap.String("chatID", chatID.String()),
		zap.Any("userIDs", userIDs))

	result, err := c.chatUsecase.DeleteUserFromChat(r.Context(), userID, userIDs, chatID)
	if err != nil {
		logger.Error("Failed to remove users from chat", zap.Error(err))
		code, err := apperrors.GetErrAndCodeToSend(err)
		if sendErr := utils.SendJSONResponse(w, code, err, false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusOK, result, true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
}

// LeaveChat позволяет пользователю покинуть чат
// @Summary Покинуть чат
// @Description Позволяет текущему пользователю покинуть чат (кроме владельца)
// @Tags Chat
// @Param chat_id path string true "ID чата"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id}/leave [post]
func (c *chatController) LeaveChat(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	vars := mux.Vars(r)
	chatID, err := uuid.Parse(vars["chat_id"])
	if err != nil {
		logger.Error("Invalid chat ID format", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid chat ID", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("User leaving chat",
		zap.String("userID", userID.String()),
		zap.String("chatID", chatID.String()))

	result, err := c.chatUsecase.DeleteUserFromChat(r.Context(), userID, []uuid.UUID{userID}, chatID)
	if err != nil {
		logger.Error("Failed to leave chat", zap.Error(err))
		code, err := apperrors.GetErrAndCodeToSend(err)
		if sendErr := utils.SendJSONResponse(w, code, err, false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusOK, result, true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
}
