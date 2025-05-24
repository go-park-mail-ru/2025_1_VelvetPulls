package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/app_errors"
	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	usecase "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/proto"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type messageController struct {
	messageUsecase usecase.IMessageUsecase
	sessionClient  authpb.SessionServiceClient
}

func NewMessageController(r *mux.Router, messageUsecase usecase.IMessageUsecase, sessionClient authpb.SessionServiceClient) {
	controller := &messageController{
		messageUsecase: messageUsecase,
		sessionClient:  sessionClient,
	}

	r.Handle("/chat/{chat_id}/messages", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.GetMessageHistory))).Methods(http.MethodGet)
	r.Handle("/chat/{chat_id}/messages/{direction}/{last_message_id}", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.GetPaginatedMessageHistory))).Methods(http.MethodGet)
	r.Handle("/chat/{chat_id}/messages", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.SendMessage))).Methods(http.MethodPost)
	r.Handle("/chat/{chat_id}/messages/{message_id}", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.UpdateMessage))).Methods(http.MethodPut)
	r.Handle("/chat/{chat_id}/messages/{message_id}", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.DeleteMessage))).Methods(http.MethodDelete)
}

// @Summary Получить историю сообщений в чате
// @Description Возвращает все сообщения в чате по chat_id
// @Tags Message
// @Produce json
// @Param chat_id path string true "ID чата"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id}/messages [get]
func (c *messageController) GetMessageHistory(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	// Парсим chat_id
	vars := mux.Vars(r)
	chatID, err := uuid.Parse(vars["chat_id"])
	if err != nil {
		logger.Error("Invalid chat ID format", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat ID", false)
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("GetMessageHistory", zap.String("chatID", chatID.String()), zap.String("userID", userID.String()))

	// Получаем историю
	messages, err := c.messageUsecase.GetChatMessages(r.Context(), userID, chatID)
	if err != nil {
		logger.Error("Failed to get message history", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, messages, true)
}

// @Summary Получить страницу сообщений в чате
// @Description Возвращает порцию сообщений чата начиная с last_message_id (исключительно), используется для пагинации
// @Tags Message
// @Produce json
// @Param chat_id path string true "ID чата"
// @Param last_message_id path string true "ID последнего сообщения на предыдущей странице"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id}/messages/{direction}/{last_message_id} [get]
func (c *messageController) GetPaginatedMessageHistory(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	vars := mux.Vars(r)

	chatID, err := uuid.Parse(vars["chat_id"])
	if err != nil {
		logger.Error("Invalid chat ID format", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat ID", false)
		return
	}

	lastMessageID, err := uuid.Parse(vars["last_message_id"])
	if err != nil {
		logger.Error("Invalid last message ID format", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid last message ID", false)
		return
	}

	direction := vars["direction"]
	if direction != "up" && direction != "down" {
		logger.Error("Invalid direction", zap.String("direction", direction))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid direction: must be 'up' or 'down'", false)
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("GetPaginatedMessageHistory",
		zap.String("direction", direction),
		zap.String("chatID", chatID.String()),
		zap.String("lastMessageID", lastMessageID.String()),
		zap.String("userID", userID.String()),
	)

	var messages []model.Message
	if direction == "up" {
		messages, err = c.messageUsecase.GetMessagesBefore(r.Context(), userID, chatID, lastMessageID)
	} else {
		messages, err = c.messageUsecase.GetMessagesAfter(r.Context(), userID, chatID, lastMessageID)
	}

	if err != nil {
		logger.Error("Failed to get paginated message history", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, messages, true)
}

// @Summary Отправить сообщение в чат
// @Description Отправляет новое сообщение в указанный чат
// @Tags Message
// @Accept multipart/form-data
// @Produce json
// @Param chat_id path string true "ID чата"
// @Param text formData string false "Текст сообщения"
// @Param sticker formData string false "Стикер (URL или ID)"
// @Param files formData file false "Файлы (можно несколько)"
// @Param photos formData file false "Фотографии (можно несколько)"
// @Success 201 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 401 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id}/messages [post]
func (c *messageController) SendMessage(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	// Парсим chat_id
	chatID, err := uuid.Parse(mux.Vars(r)["chat_id"])
	if err != nil {
		logger.Error("Invalid chat ID format", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat ID", false)
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())
	logger.Info("SendMessage", zap.String("chatID", chatID.String()), zap.String("userID", userID.String()))

	// Парсим multipart/form
	if err := r.ParseMultipartForm(config.MAX_FILE_SIZE); err != nil {
		logger.Error("Failed to parse multipart form", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Request too large or malformed", false)
		return
	}

	var input model.MessageInput
	// Читаем текст и стикер
	if data := r.FormValue("text"); data != "" {
		if err := json.Unmarshal([]byte(data), &input); err != nil {
			logger.Error("Invalid chat data format", zap.Error(err))
			utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat data format", false)
			return
		}
	}
	sticker := r.FormValue("sticker")

	// Логируем полученные данные
	logger.Info("Parsed message form values", zap.String("sticker", sticker))

	var msg model.Message
	msg.Body = input.Message
	msg.ChatID = chatID
	msg.UserID = userID
	msg.Sticker = sticker

	files := r.MultipartForm.File["files"]
	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			logger.Error("Failed to open file", zap.Error(err))
			utils.SendJSONResponse(w, r, http.StatusBadRequest, "Failed to open files", false)
			return
		}
		defer file.Close()

		msg.Files = append(msg.Files, file)
		msg.FilesHeaders = append(msg.FilesHeaders, header)
	}

	photos := r.MultipartForm.File["photos"]
	for _, header := range photos {
		photo, err := header.Open()
		if err != nil {
			logger.Error("Failed to open photo", zap.Error(err))
			utils.SendJSONResponse(w, r, http.StatusBadRequest, "Failed to open photos", false)
			return
		}
		defer photo.Close()

		msg.Photos = append(msg.Photos, photo)
		msg.PhotosHeaders = append(msg.PhotosHeaders, header)
	}

	// Вызываем usecase
	if err := c.messageUsecase.SendMessage(r.Context(), &msg, userID, chatID); err != nil {
		logger.Error("Failed to send message", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusCreated, "Message sent successfully", true)
}

// @Summary Обновить сообщение в чате
// @Description Обновляет сообщение пользователя в чате
// @Tags Message
// @Accept json
// @Produce json
// @Param chat_id path string true "ID чата"
// @Param message body model.MessageInput true "Новое сообщение"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id}/messages [put]
func (c *messageController) UpdateMessage(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	vars := mux.Vars(r)
	chatID, err := uuid.Parse(vars["chat_id"])
	if err != nil {
		logger.Error("Invalid chat ID format", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat ID", false)
		return
	}

	messageID, err := uuid.Parse(vars["message_id"])
	if err != nil {
		logger.Error("Invalid message ID format", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid message ID", false)
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())

	var input model.MessageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("Failed to decode message input", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid request body", false)
		return
	}

	logger.Info("UpdateMessage", zap.String("chatID", chatID.String()), zap.String("userID", userID.String()))

	if err := c.messageUsecase.UpdateMessage(r.Context(), messageID, &input, userID, chatID); err != nil {
		logger.Error("Failed to update message", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, "Message updated successfully", true)
}

// @Summary Удалить сообщение в чате
// @Description Удаляет сообщение пользователя из чата
// @Tags Message
// @Produce json
// @Param chat_id path string true "ID чата"
// @Param message_id path string true "ID сообщения"
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 403 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id}/messages/{message_id} [delete]
func (c *messageController) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	vars := mux.Vars(r)
	chatID, err := uuid.Parse(vars["chat_id"])
	if err != nil {
		logger.Error("Invalid chat ID format", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid chat ID", false)
		return
	}

	messageID, err := uuid.Parse(vars["message_id"])
	if err != nil {
		logger.Error("Invalid message ID format", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid message ID", false)
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())

	logger.Info("DeleteMessage", zap.String("chatID", chatID.String()), zap.String("userID", userID.String()), zap.String("messageID", messageID.String()))

	if err := c.messageUsecase.DeleteMessage(r.Context(), messageID, userID, chatID); err != nil {
		logger.Error("Failed to delete message", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		utils.SendJSONResponse(w, r, code, msg, false)
		return
	}

	utils.SendJSONResponse(w, r, http.StatusOK, "Message deleted successfully", true)
}
