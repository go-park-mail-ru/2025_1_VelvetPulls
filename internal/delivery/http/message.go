package http

import (
	"encoding/json"
	"net/http"

	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/app_errors"
	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	usecase "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type messageController struct {
	messageUsecase usecase.IMessageUsecase
	sessionUsecase usecase.ISessionUsecase
}

func NewMessageController(r *mux.Router, messageUsecase usecase.IMessageUsecase, sessionUsecase usecase.ISessionUsecase) {
	controller := &messageController{
		messageUsecase: messageUsecase,
		sessionUsecase: sessionUsecase,
	}

	r.Handle("/chat/{chat_id}/messages", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.GetMessageHistory))).Methods(http.MethodGet)
	r.Handle("/chat/{chat_id}/messages", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.SendMessage))).Methods(http.MethodPost)
	r.Handle("/chat/{chat_id}/messages/{message_id}", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.UpdateMessage))).Methods(http.MethodPut)
	r.Handle("/chat/{chat_id}/messages/{message_id}", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.DeleteMessage))).Methods(http.MethodDelete)
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

// @Summary Отправить сообщение в чат
// @Description Отправляет новое сообщение в указанный чат
// @Tags Message
// @Accept json
// @Produce json
// @Param chat_id path string true "ID чата"
// @Param message body model.MessageInput true "Сообщение"
// @Success 201 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 401 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /chat/{chat_id}/messages [post]
func (c *messageController) SendMessage(w http.ResponseWriter, r *http.Request) {
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

	// Декодируем тело
	var input model.MessageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("Failed to decode message input", zap.Error(err))
		utils.SendJSONResponse(w, r, http.StatusBadRequest, "Invalid request body", false)
		return
	}

	logger.Info("SendMessage", zap.String("chatID", chatID.String()), zap.String("userID", userID.String()))

	// Отправляем через usecase
	if err := c.messageUsecase.SendMessage(r.Context(), &input, userID, chatID); err != nil {
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
