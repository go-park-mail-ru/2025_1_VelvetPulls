package repository

import (
	"slices"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
)

var chats = []model.Chat{
	{
		ID:          1,
		Type:        model.ChatTypeDialog,
		Title:       "Chat 1",
		Description: "Description of chat 1: user 1 and 2",
		Members:     []int64{1, 2},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          2,
		Type:        model.ChatTypeDialog,
		Title:       "Chat 2",
		Description: "Description of chat 2: users 2 and 3",
		Members:     []int64{2, 3},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          3,
		Type:        model.ChatTypeDialog,
		Title:       "Chat 3",
		Description: "Description of chat 3: users 1 and 3",
		Members:     []int64{1, 3},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}

func GetChatByID(chatid int64) (model.Chat, error) {
	for _, chat := range chats {
		if chat.ID == chatid {
			return chat, nil
		}
	}

	return model.Chat{}, apperrors.ErrChatNotFound
}

func GetChatsWithUser(userid int64) []model.Chat {
	result := make([]model.Chat, 0)

	for _, chat := range chats {
		if slices.Contains(chat.Members, userid) {
			result = append(result, chat)
		}
	}

	return result
}
