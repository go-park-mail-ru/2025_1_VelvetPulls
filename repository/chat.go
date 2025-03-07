package repository

import (
	"slices"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/models"
)

var chats = []models.Chat{
	{
		ID:          1,
		Type:        models.ChatTypeDialog,
		Title:       "Chat 1",
		Description: "Description of chat 1: user 1 and 2",
		Members:     []int64{1, 2},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          2,
		Type:        models.ChatTypeDialog,
		Title:       "Chat 2",
		Description: "Description of chat 2: users 2 and 3",
		Members:     []int64{2, 3},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          3,
		Type:        models.ChatTypeDialog,
		Title:       "Chat 3",
		Description: "Description of chat 3: users 1 and 3",
		Members:     []int64{1, 3},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}

func GetChatByID(chatid int64) (models.Chat, error) {
	for _, chat := range chats {
		if chat.ID == chatid {
			return chat, nil
		}
	}

	return models.Chat{}, apperrors.ErrChatNotFound
}

func GetChatsWithUser(userid int64) []models.Chat {
	result := make([]models.Chat, 0)

	for _, chat := range chats {
		if slices.Contains(chat.Members, userid) {
			result = append(result, chat)
		}
	}

	return result
}
