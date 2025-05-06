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
) (map[string][]model.Chat, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, model.ErrInvalidUUID
	}

	chats, participants, err := uc.chatRepo.SearchUserChats(ctx, userID, query)
	if err != nil {
		return nil, err
	}

	result := map[string][]model.Chat{
		"dialogs":  make([]model.Chat, 0),
		"groups":   make([]model.Chat, 0),
		"channels": make([]model.Chat, 0),
	}

	for i, chat := range chats {
		if chat.Type == "dialog" {
			uc.decorateDialogInfo(&chat, userID, participants[i])
		}

		// Всегда добавляем во все категории
		switch chat.Type {
		case "dialog":
			result["dialogs"] = append(result["dialogs"], chat)
		case "group":
			result["groups"] = append(result["groups"], chat)
		case "channel":
			result["channels"] = append(result["channels"], chat)
		}
	}

	metrics.IncBusinessOp("search_chats")
	return result, nil
}

func (uc *ChatUsecase) decorateDialogInfo(chat *model.Chat, me uuid.UUID, users []model.UserInChat) {
	otherUsers := make([]model.UserInChat, 0)
	for _, u := range users {
		if u.ID != me {
			otherUsers = append(otherUsers, u)
		}
	}

	switch len(otherUsers) {
	case 0:
		chat.Title = "Saved Messages"
		chat.AvatarPath = nil
	case 1:
		chat.Title = otherUsers[0].Username
		chat.AvatarPath = otherUsers[0].AvatarPath
	default:
		chat.Title = "lolChto?"
		chat.AvatarPath = nil
	}
}
