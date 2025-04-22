package usecase

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// IMessageUsecase описывает операции по работе с сообщениями
// в рамках чата: чтение истории и отправка новых сообщений
type IMessageUsecase interface {
	GetChatMessages(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) ([]model.Message, error)
	SendMessage(ctx context.Context, input *model.MessageInput, userID uuid.UUID, chatID uuid.UUID) error
}

// MessageUsecase реализует бизнес-логику сообщений в чатах
type MessageUsecase struct {
	messageRepo repository.IMessageRepo
	chatRepo    repository.IChatRepo
	wsUsecase   IWebsocketUsecase
}

// NewMessageUsecase создаёт экземпляр MessageUsecase
func NewMessageUsecase(msgRepo repository.IMessageRepo, chatRepo repository.IChatRepo, wsUsecase IWebsocketUsecase) IMessageUsecase {
	return &MessageUsecase{messageRepo: msgRepo, chatRepo: chatRepo, wsUsecase: wsUsecase}
}

// ensureMember проверяет, что пользователь является участником чата
func (uc *MessageUsecase) ensureMember(ctx context.Context, userID, chatID uuid.UUID) error {
	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		return err
	}
	if !model.UserRoleInChat(role).IsMember() {
		return ErrPermissionDenied
	}
	return nil
}

// GetChatMessages возвращает все сообщения из чата
func (uc *MessageUsecase) GetChatMessages(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) ([]model.Message, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("GetChatMessages start", zap.String("userID", userID.String()), zap.String("chatID", chatID.String()))

	// проверяем права доступа
	if err := uc.ensureMember(ctx, userID, chatID); err != nil {
		logger.Warn("Access denied при попытке получить сообщения", zap.Error(err))
		return nil, err
	}

	// выборка из репозитория
	msgs, err := uc.messageRepo.GetMessages(ctx, chatID)
	if err != nil {
		logger.Error("GetMessages failed", zap.Error(err))
		return nil, err
	}
	return msgs, nil
}

// SendMessage валидирует и сохраняет новое сообщение, а затем публикует событие
func (uc *MessageUsecase) SendMessage(ctx context.Context, input *model.MessageInput, userID uuid.UUID, chatID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("SendMessage start", zap.String("userID", userID.String()), zap.String("chatID", chatID.String()))

	// проверяем права доступа
	if err := uc.ensureMember(ctx, userID, chatID); err != nil {
		logger.Warn("Access denied при попытке отправить сообщение", zap.Error(err))
		return err
	}

	// валидация содержимого
	if err := input.Validate(); err != nil {
		logger.Error("Validation failed для MessageInput", zap.Error(err))
		return fmt.Errorf("SendMessage: validation failed: %w", err)
	}

	// подготовка модели сообщения
	msg := &model.Message{
		ChatID: chatID,
		UserID: userID,
		Body:   input.Message,
	}

	// сохранение через репозиторий
	saved, err := uc.messageRepo.CreateMessage(ctx, msg)
	if err != nil {
		logger.Error("CreateMessage failed", zap.Error(err))
		return fmt.Errorf("SendMessage: failed to create message: %w", err)
	}

	// публикация события WebSocket/NATS
	uc.publishEvent(ctx, NewMessage, saved)
	return nil
}

// publishEvent формирует и отправляет событие о новом сообщении
func (uc *MessageUsecase) publishEvent(ctx context.Context, action string, message *model.Message) {
	logger := utils.GetLoggerFromCtx(ctx)
	event := model.MessageEvent{Action: action, Message: *message}
	if uc.wsUsecase == nil {
		logger.Warn("wsUsecase не инициализирован, событие не отправлено")
		return
	}
	if err := uc.wsUsecase.PublishMessage(event); err != nil {
		logger.Error("PublishMessage failed", zap.Error(err))
	} else {
		logger.Info("Message event published", zap.String("chatID", message.ChatID.String()))
	}
}
