package usecase

import (
	"context"
	"sync"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/google/uuid"
)

const (
	AddWebsocketUser    = "addWebsocketUser"
	DeleteWebsocketUser = "deleteWebsocketUser"
)

type IWebsocketUsecase interface {
	RegisterUserChannel(userID uuid.UUID, eventChan chan model.AnyEvent) error
	UnregisterUserChannel(userID uuid.UUID, eventChan chan model.AnyEvent)
	GetUserChannels(userID uuid.UUID) []chan model.AnyEvent
	InitChat(chatID uuid.UUID, users []uuid.UUID) error
	ConsumeMessages()
	ConsumeChats()
	SendMessage(event model.MessageEvent)
	SendChatEvent(event model.ChatEvent)
}

type WebsocketUsecase struct {
	messageChan chan model.MessageEvent
	chatChan    chan model.ChatEvent

	onlineChats map[uuid.UUID]model.ChatInfoWS
	onlineUsers map[uuid.UUID][]chan model.AnyEvent
	mu          sync.RWMutex

	chatRepo repository.IChatRepo
}

func NewWebsocketUsecase(chatRepo repository.IChatRepo) IWebsocketUsecase {
	return &WebsocketUsecase{
		messageChan: make(chan model.MessageEvent, 100),
		chatChan:    make(chan model.ChatEvent, 100),
		onlineChats: make(map[uuid.UUID]model.ChatInfoWS),
		onlineUsers: make(map[uuid.UUID][]chan model.AnyEvent),
		chatRepo:    chatRepo,
	}
}

func (w *WebsocketUsecase) GetUserChannels(userID uuid.UUID) []chan model.AnyEvent {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.onlineUsers[userID]
}

func (w *WebsocketUsecase) RegisterUserChannel(userID uuid.UUID, eventChan chan model.AnyEvent) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Добавляем новый канал в слайс
	w.onlineUsers[userID] = append(w.onlineUsers[userID], eventChan)

	// Получаем все чаты пользователя
	userChats, _, err := w.chatRepo.GetChats(context.Background(), userID)
	if err != nil {
		return err
	}

	for _, chat := range userChats {
		if _, ok := w.onlineChats[chat.ID]; !ok {
			w.initNewChatRoom(chat.ID)
		}
		w.onlineChats[chat.ID].Events <- model.Event{
			Action: AddWebsocketUser,
			ChatId: chat.ID,
			Users:  []uuid.UUID{userID},
		}
	}

	return nil
}

func (w *WebsocketUsecase) UnregisterUserChannel(userID uuid.UUID, eventChan chan model.AnyEvent) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Удаляем канал из списка пользователя
	channels := w.onlineUsers[userID]
	for i, ch := range channels {
		if ch == eventChan {
			w.onlineUsers[userID] = append(channels[:i], channels[i+1:]...)
			break
		}
	}

	// Если у пользователя больше нет каналов, удаляем его из чатов
	if len(w.onlineUsers[userID]) == 0 {
		delete(w.onlineUsers, userID)

		// Удаляем пользователя из всех чатов
		for chatID, chatInfo := range w.onlineChats {
			if _, ok := chatInfo.Users[userID]; ok {
				delete(chatInfo.Users, userID)

				// Если в чате больше нет пользователей, закрываем его
				if len(chatInfo.Users) == 0 {
					close(chatInfo.Events)
					delete(w.onlineChats, chatID)
				}
			}
		}
	}

	close(eventChan)
}

func (w *WebsocketUsecase) initNewChatRoom(chatID uuid.UUID) {
	eventsChan := make(chan model.Event, 100)
	chatInfo := model.ChatInfoWS{
		Events: eventsChan,
		Users:  make(map[uuid.UUID]struct{}),
	}
	w.onlineChats[chatID] = chatInfo

	go func(chatInfo model.ChatInfoWS) {
		for event := range chatInfo.Events {
			switch event.Action {
			case AddWebsocketUser:
				for _, userID := range event.Users {
					chatInfo.Users[userID] = struct{}{}
				}
			case DeleteWebsocketUser:
				for _, userID := range event.Users {
					delete(chatInfo.Users, userID)
				}
			}
		}
	}(chatInfo)
}

func (w *WebsocketUsecase) InitChat(chatID uuid.UUID, users []uuid.UUID) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, exists := w.onlineChats[chatID]; !exists {
		w.initNewChatRoom(chatID)
	}

	for _, userID := range users {
		if _, exists := w.onlineUsers[userID]; exists {
			w.onlineChats[chatID].Events <- model.Event{
				Action: AddWebsocketUser,
				Users:  []uuid.UUID{userID},
				ChatId: chatID,
			}
		}
	}

	return nil
}
