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

type IMessageUsecase interface {
	GetChatMessages(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) ([]model.Message, error)
	SendMessage(ctx context.Context, messageInput *model.MessageInput, userID uuid.UUID, chatID uuid.UUID) error
}

type MessageUsecase struct {
	messageRepo repository.IMessageRepo
	chatRepo    repository.IChatRepo
}

func NewMessageUsecase(messageRepo repository.IMessageRepo, chatRepo repository.IChatRepo) IMessageUsecase {
	return &MessageUsecase{messageRepo: messageRepo, chatRepo: chatRepo}
}

func (uc *MessageUsecase) GetChatMessages(ctx context.Context, userID uuid.UUID, chatID uuid.UUID) ([]model.Message, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		logger.Error("Failed to get user role in chat", zap.Error(err))
		return nil, err
	}
	if role != "owner" && role != "member" {
		logger.Warn("Permission denied for user in chat",
			zap.String("role", role))
		return nil, ErrPermissionDenied
	}

	logger.Info("Fetching messages for chat: " + chatID.String())
	return uc.messageRepo.GetMessages(ctx, chatID)
}

func (uc *MessageUsecase) SendMessage(ctx context.Context, messageInput *model.MessageInput, userID uuid.UUID, chatID uuid.UUID) error {
	logger := utils.GetLoggerFromCtx(ctx)

	role, err := uc.chatRepo.GetUserRoleInChat(ctx, userID, chatID)
	if err != nil {
		logger.Error("Failed to get user role in chat", zap.Error(err))
		return err
	}
	if role != "owner" && role != "member" {
		logger.Warn("Permission denied for user in chat",
			zap.String("role", role))
		return ErrPermissionDenied
	}

	if err := messageInput.Validate(); err != nil {
		logger.Error("Invalid message payload", zap.Error(err))
		return fmt.Errorf("SendMessage: validation failed: %w", err)
	}

	message := &model.Message{
		ChatID: chatID,
		UserID: userID,
		Body:   messageInput.Message,
	}

	if err := uc.messageRepo.CreateMessage(ctx, message); err != nil {
		logger.Error("Failed to create message", zap.Error(err))
		return fmt.Errorf("SendMessage: failed to create message: %w", err)
	}
	return nil
}
