package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/repository"
	"github.com/google/uuid"
)

type ChatUsecase struct {
	chatRepo repository.ChatRepo
}

func NewChatUsecase(chatRepo repository.ChatRepo) *ChatUsecase {
	return &ChatUsecase{chatRepo: chatRepo}
}

func (uc *ChatUsecase) SearchUserChats(
	ctx context.Context,
	userIDStr string,
	query string,
	types []string,
) ([]model.Chat, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, model.ErrValidation
	}

	validTypes := make([]string, 0, len(types))
	for _, t := range types {
		if !isValidChatType(t) {
			return nil, ErrChatType
		}
		validTypes = append(validTypes, t)
	}

	metrics.IncBusinessOp("search_chats")
	return uc.chatRepo.SearchUserChats(ctx, userID, query, validTypes)
}

func isValidChatType(t string) bool {
	switch t {
	case "dialog", "group", "channel":
		return true
	default:
		return false
	}
}
