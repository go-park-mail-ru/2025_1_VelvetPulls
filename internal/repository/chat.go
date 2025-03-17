package repository

import (
	"sync"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
)

type ChatRepoInterface interface {
	GetChatsByUsername(username string) ([]model.Chat, error)
	AddChat(chat *model.Chat)
}

type сhatRepo struct {
	chats []*model.Chat
	mu    sync.RWMutex // Мьютекс для безопасного чтения и записи
}

func NewChatRepo() ChatRepoInterface {
	return &сhatRepo{
		chats: make([]*model.Chat, 0),
	}
}

// Получение чатов по имени пользователя (возвращает копии чатов, безопасно для конкурентного чтения)
func (r *сhatRepo) GetChatsByUsername(username string) ([]model.Chat, error) {
	r.mu.RLock() // Блокируем только для чтения
	defer r.mu.RUnlock()

	userChats := make([]model.Chat, 0)
	for _, chat := range r.chats {
		if chat.OwnerUsername == username {
			userChats = append(userChats, *chat) // Возвращаем копию чата
		}
	}

	return userChats, nil
}

// Добавление нового чата (безопасно для конкурентной записи)
func (r *сhatRepo) AddChat(chat *model.Chat) {
	r.mu.Lock() // Блокируем для записи
	defer r.mu.Unlock()

	chat.CreatedAt = time.Now()
	chat.UpdatedAt = chat.CreatedAt
	r.chats = append(r.chats, chat)
}
