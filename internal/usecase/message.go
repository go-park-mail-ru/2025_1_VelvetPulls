package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type IMessageUsecase interface {
	GetChatMessages(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) ([]model.Message, error)
	GetMessagesBefore(ctx context.Context, userID, chatID, messageID uuid.UUID) ([]model.Message, error)
	SendMessage(ctx context.Context, input *model.MessageInput, userID uuid.UUID, chatID uuid.UUID) error
	UpdateMessage(ctx context.Context, messageID uuid.UUID, input *model.MessageInput, userID uuid.UUID, chatID uuid.UUID) error
	DeleteMessage(ctx context.Context, messageID uuid.UUID, userID uuid.UUID, chatID uuid.UUID) error
}

type MessageUsecase struct {
	messageRepo repository.IMessageRepo
	chatRepo    repository.IChatRepo
	nc          *nats.Conn
}

func NewMessageUsecase(msgRepo repository.IMessageRepo, chatRepo repository.IChatRepo, nc *nats.Conn) IMessageUsecase {
	return &MessageUsecase{messageRepo: msgRepo, chatRepo: chatRepo, nc: nc}
}

func (uc *MessageUsecase) GetChatMessages(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) ([]model.Message, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("GetChatMessages start", zap.String("userID", userID.String()), zap.String("chatID", chatID.String()))

	if err := uc.ensureMember(ctx, userID, chatID); err != nil {
		logger.Warn("Access denied при попытке получить сообщения", zap.Error(err))
		return nil, err
	}

	msgs, err := uc.messageRepo.GetMessages(ctx, chatID)
	if err != nil {
		logger.Error("GetMessages failed", zap.Error(err))
		return nil, err
	}
	metrics.IncBusinessOp("get_messages")
	return msgs, nil
}

func (uc *MessageUsecase) GetMessagesBefore(ctx context.Context, userID, chatID, messageID uuid.UUID) ([]model.Message, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("GetMessagesAfter start", zap.String("userID", userID.String()), zap.String("chatID", chatID.String()), zap.String("messageID", messageID.String()))

	if err := uc.ensureMember(ctx, userID, chatID); err != nil {
		logger.Warn("Access denied при попытке получить сообщения до", zap.Error(err))
		return nil, err
	}

	messages, err := uc.messageRepo.GetMessagesBefore(ctx, chatID, messageID)
	if err != nil {
		logger.Error("GetMessagesAfterID failed", zap.Error(err))
		return nil, err
	}

	metrics.IncBusinessOp("get_messages_before")
	return messages, nil
}

func (uc *MessageUsecase) SendMessage(ctx context.Context, input *model.MessageInput, userID uuid.UUID, chatID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("SendMessage start", zap.String("userID", userID.String()), zap.String("chatID", chatID.String()))

	if err := uc.ensureCanSend(ctx, userID, chatID); err != nil {
		logger.Warn("Access denied при попытке отправить сообщение", zap.Error(err))
		return err
	}

	if err := input.Validate(); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrMessageValidationFailed, err)
	}

	msg := &model.Message{ChatID: chatID, UserID: userID, Body: input.Message}

	saved, err := uc.messageRepo.CreateMessage(ctx, msg)
	if err != nil {
		logger.Error("CreateMessage failed", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrMessageCreationFailed, err)
	}

	// publish new-message event
	e := model.MessageEvent{Action: utils.NewMessage, Message: *saved}
	data, _ := json.Marshal(e)
	subj := fmt.Sprintf("chat.%s.messages", chatID.String())
	if err := uc.nc.Publish(subj, data); err != nil {
		logger.Error("NATS publish failed", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrMessagePublishFailed, err)
	}

	metrics.IncBusinessOp("send_message")
	return nil
}

func (uc *MessageUsecase) UpdateMessage(ctx context.Context, messageID uuid.UUID, input *model.MessageInput, userID uuid.UUID, chatID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("UpdateMessage start", zap.String("userID", userID.String()), zap.String("chatID", chatID.String()))

	if err := uc.ensureMember(ctx, userID, chatID); err != nil {
		logger.Warn("Access denied при попытке редактировать сообщение", zap.Error(err))
		return err
	}

	if err := input.Validate(); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrMessageValidationFailed, err)
	}

	message, err := uc.messageRepo.GetMessage(ctx, messageID)
	if err != nil {
		logger.Error("GetMessage failed", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrMessageNotFound, err)
	}

	if message.UserID != userID {
		logger.Warn("Access denied: user is not the author of the message", zap.String("messageID", messageID.String()))
		return ErrMessageAccessDenied
	}

	updated, err := uc.messageRepo.UpdateMessage(ctx, messageID, input.Message)
	if err != nil {
		logger.Error("UpdateMessage failed", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrMessageUpdateFailed, err)
	}

	// publish update-message event
	e := model.MessageEvent{Action: utils.UpdateMessage, Message: *updated}
	data, _ := json.Marshal(e)
	subj := fmt.Sprintf("chat.%s.messages", chatID.String())
	if err := uc.nc.Publish(subj, data); err != nil {
		logger.Error("NATS publish failed", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrMessagePublishFailed, err)
	}

	metrics.IncBusinessOp("update_message")
	return nil
}

func (uc *MessageUsecase) DeleteMessage(ctx context.Context, messageID uuid.UUID, userID uuid.UUID, chatID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("DeleteMessage start", zap.String("userID", userID.String()), zap.String("chatID", chatID.String()))

	if err := uc.ensureMember(ctx, userID, chatID); err != nil {
		logger.Warn("Access denied при попытке удалить сообщение", zap.Error(err))
		return err
	}

	message, err := uc.messageRepo.GetMessage(ctx, messageID)
	if err != nil {
		logger.Error("GetMessage failed", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrMessageNotFound, err)
	}

	if message.UserID != userID {
		logger.Warn("Access denied: user is not the author of the message", zap.String("messageID", messageID.String()))
		return ErrMessageAccessDenied
	}

	deleted, err := uc.messageRepo.DeleteMessage(ctx, messageID)
	if err != nil {
		logger.Error("DeleteMessage failed", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrMessageDeleteFailed, err)
	}

	// publish delete-message event
	e := model.MessageEvent{Action: utils.DeleteMessage, Message: *deleted}
	data, _ := json.Marshal(e)
	subj := fmt.Sprintf("chat.%s.messages", chatID.String())
	if err := uc.nc.Publish(subj, data); err != nil {
		logger.Error("NATS publish failed", zap.Error(err))
		return fmt.Errorf("%w: %v", ErrMessagePublishFailed, err)
	}

	metrics.IncBusinessOp("delete_message")
	return nil
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

// ensureCanSend проверяет, что пользователь может отправлять сообщения
func (uc *MessageUsecase) ensureCanSend(ctx context.Context, userID, chatID uuid.UUID) error {
	chat, err := uc.chatRepo.GetChatByID(ctx, chatID)
	if err != nil {
		return err
	}
	if chat.Type == string(model.ChatTypeChannel) {
		role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
		if err != nil {
			return err
		}
		if model.UserRoleInChat(role) != model.RoleOwner {
			return ErrMessageAccessDenied
		}
		return nil
	}
	return uc.ensureMember(ctx, userID, chatID)
}
