package http

import (
	"encoding/json"
	"net/http"

	apperrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/app_errors"
	model "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	usecase "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/delivery/proto"
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
	r.Handle("/chat/{chat_id}/messages", middleware.AuthMiddleware(sessionClient)(http.HandlerFunc(controller.SendMessage))).Methods(http.MethodPost)
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
	logger.Info("Fetching message history", zap.String("chatID", chatID.String()), zap.String("userID", userID.String()))

	messages, err := c.messageUsecase.GetChatMessages(r.Context(), userID, chatID)
	if err != nil {
		logger.Error("Failed to get message history", zap.Error(err))
		code, msg := apperrors.GetErrAndCodeToSend(err)
		if sendErr := utils.SendJSONResponse(w, code, msg, false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusOK, messages, true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
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
	userID := utils.GetUserIDFromCtx(r.Context())

	vars := mux.Vars(r)
	chatID, err := uuid.Parse(vars["chat_id"])
	if err != nil {
		logger.Error("Invalid chat ID format", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid chat ID", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	var input model.MessageInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Error("Failed to decode message input", zap.Error(err))
		if sendErr := utils.SendJSONResponse(w, http.StatusBadRequest, "Invalid request body", false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	logger.Info("Sending message", zap.String("userID", userID.String()), zap.String("chatID", chatID.String()))

	err = c.messageUsecase.SendMessage(r.Context(), &input, userID, chatID)
	if err != nil {
		logger.Error("Failed to send message", zap.Error(err))
		code, errMsg := apperrors.GetErrAndCodeToSend(err)
		if sendErr := utils.SendJSONResponse(w, code, errMsg, false); sendErr != nil {
			logger.Error("Failed to send error response", zap.Error(sendErr))
		}
		return
	}

	if sendErr := utils.SendJSONResponse(w, http.StatusCreated, "message send successful", true); sendErr != nil {
		logger.Error("Failed to send success response", zap.Error(sendErr))
	}
}
