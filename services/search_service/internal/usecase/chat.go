package usecase

import (
	"context"
	"strings"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ChatUsecase struct {
	chatRepo repository.ChatRepo
}

func NewChatUsecase(chatRepo repository.ChatRepo) *ChatUsecase {
	return &ChatUsecase{chatRepo: chatRepo}
}

func filterChats(chats []model.Chat, query string) []model.Chat {
	filtered := make([]model.Chat, 0)
	query = strings.ToLower(query)

	for _, c := range chats {
		if strings.Contains(strings.ToLower(c.Title), query) {
			filtered = append(filtered, c)
		}
	}

	return filtered
}

func (uc *ChatUsecase) decorateDialog(ctx context.Context, chat *model.Chat, me uuid.UUID) {
	users, err := uc.chatRepo.GetUsersFromChat(ctx, chat.ID)
	if err != nil {
		zap.L().Warn("decorateDialog: GetUsersFromChat failed", zap.Error(err))
		return
	}

	if len(users) == 1 {
		chat.Title = users[0].Username
		chat.AvatarPath = users[0].AvatarPath
		return
	}

	for _, u := range users {
		if u.ID != me {
			chat.Title = u.Username
			chat.AvatarPath = u.AvatarPath
			break
		}
	}
}

func (uc *ChatUsecase) SearchUserChats(
	ctx context.Context,
	userIDStr string,
	query string,
) (*model.ChatGroups, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, model.ErrValidation
	}

	chats, err := uc.chatRepo.SearchUserChats(ctx, userID, query)
	if err != nil {
		return nil, err
	}

	groups := &model.ChatGroups{
		Dialogs:  []model.Chat{},
		Groups:   []model.Chat{},
		Channels: []model.Chat{},
	}

	for _, chat := range chats {
		switch chat.Type {
		case "group":
			groups.Groups = append(groups.Groups, chat)
		case "channel":
			groups.Channels = append(groups.Channels, chat)
		}
	}

	allChats, err := uc.chatRepo.SearchUserChats(ctx, userID, "")
	if err != nil {
		return nil, err
	}

	for _, chat := range allChats {
		if chat.Type == "dialog" {
			uc.decorateDialog(ctx, &chat, userID)
			groups.Dialogs = append(groups.Dialogs, chat)
		}
	}

	groups.Dialogs = filterChats(groups.Dialogs, query)

	metrics.IncBusinessOp("search_chats")
	return groups, nil
}
