package repository

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
)

var (
	chats = []*model.Chat{
		{
			OwnerUsername: "ruslantus228",
			Type:          model.ChatTypeDialog,
			Title:         "Chat 1",
			Description:   "Description of chat 1: user 1 and 2",
			Members:       []int64{1, 2},
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OwnerUsername: "ruslantus228",
			Type:          model.ChatTypeDialog,
			Title:         "Chat 2",
			Description:   "Description of chat 2: users 2 and 3",
			Members:       []int64{2, 3},
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			OwnerUsername: "ruslantus228",
			Type:          model.ChatTypeDialog,
			Title:         "Chat 3",
			Description:   "Description of chat 3: users 1 and 3",
			Members:       []int64{1, 3},
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}
	muChats sync.RWMutex // RWMutex для чтения и записи
)

// Получение чатов по имени пользователя (возвращает указатели на чаты, безопасно для конкурентного чтения)
func GetChatsByUsername(username string) ([]*model.Chat, error) {
	muChats.RLock() // Блокируем только для чтения
	defer muChats.RUnlock()

	var userChats = make([]*model.Chat, 0)
	for _, chat := range chats {
		if chat.OwnerUsername == username {
			userChats = append(userChats, chat)
		}
	}

	return userChats, nil
}

// Добавление нового чата (безопасно для конкурентной записи)
func AddChat(chat *model.Chat) {
	muChats.Lock() // Блокируем для записи
	defer muChats.Unlock()

	chat.CreatedAt = time.Now()
	chat.UpdatedAt = chat.CreatedAt
	chats = append(chats, chat)
}
