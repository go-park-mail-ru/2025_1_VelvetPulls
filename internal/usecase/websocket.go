package usecase

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/google/uuid"
)

type MessageEvent struct {
	Action  string        `json:"action"`
	Message model.Message `json:"payload"`
}

const (
	NewMessage       = "newMessage"
	AddWebsocketUser = "addWebsocketUser"
)

type AnyEvent struct {
	TypeOfEvent string
	Event       interface{}
}

type Event struct {
	Action string      `json:"action"`
	ChatId uuid.UUID   `json:"chatId"`
	Users  []uuid.UUID `json:"users"`
}

type ChatInfo struct {
	events chan Event
	users  map[uuid.UUID]struct{}
}

func SerializeMessageEvent(event MessageEvent) ([]byte, error) {
	return json.Marshal(event)
}

func DeserializeMessageEvent(data []byte) (MessageEvent, error) {
	var event MessageEvent
	err := json.Unmarshal(data, &event)
	if err != nil {
		return MessageEvent{}, err
	}
	return event, nil
}

type IWebsocketUsecase interface {
	InitBrokersForUser(userID uuid.UUID, eventChan chan AnyEvent) error
	ConsumeMessages()
	SendMessage(event MessageEvent)
}

type WebsocketUsecase struct {
	messageChan chan MessageEvent

	onlineChats map[uuid.UUID]ChatInfo
	onlineUsers map[uuid.UUID]chan AnyEvent

	chatRepo repository.IChatRepo
}

func NewWebsocketUsecase(chatRepo repository.IChatRepo) IWebsocketUsecase {
	socket := &WebsocketUsecase{
		messageChan: make(chan MessageEvent, 100),
		onlineChats: make(map[uuid.UUID]ChatInfo),
		onlineUsers: make(map[uuid.UUID]chan AnyEvent),
		chatRepo:    chatRepo,
	}

	return socket
}

func (w *WebsocketUsecase) InitNewChatRoom(chatID uuid.UUID) {
	log.Printf("Initializing new chat room: %s", chatID.String())

	eventsChan := make(chan Event, 100)
	chatInfo := ChatInfo{
		events: eventsChan,
		users:  make(map[uuid.UUID]struct{}),
	}
	w.onlineChats[chatID] = chatInfo

	go func(chatInfo ChatInfo) {
		for event := range chatInfo.events {
			log.Printf("Chat %s received event: %v", chatID.String(), event)

			switch event.Action {
			case AddWebsocketUser:
				for _, userID := range event.Users {
					chatInfo.users[userID] = struct{}{}
					log.Printf("User %s added to chat %s", userID.String(), chatID.String())
				}
			}
		}
	}(chatInfo)
}

func (w *WebsocketUsecase) InitBrokersForUser(userID uuid.UUID, eventChan chan AnyEvent) error {
	// Сохраняем канал пользователя
	w.onlineUsers[userID] = eventChan

	// Получаем все чаты пользователя
	userChats, _, err := w.chatRepo.GetChats(context.Background(), userID)
	if err != nil {
		return err
	}

	// Инициализируем комнату, если еще нет, и добавляем пользователя
	for _, chat := range userChats {
		if _, ok := w.onlineChats[chat.ID]; !ok {
			w.InitNewChatRoom(chat.ID)
		}
		w.onlineChats[chat.ID].events <- Event{
			Action: AddWebsocketUser,
			ChatId: chat.ID,
			Users:  []uuid.UUID{userID},
		}
	}

	return nil
}

func (w *WebsocketUsecase) ConsumeMessages() {
	for msg := range w.messageChan {
		chatID := msg.Message.ChatID

		if _, ok := w.onlineChats[chatID]; !ok {
			w.InitNewChatRoom(chatID)
		}

		w.SendMessage(msg)
	}
}

func (w *WebsocketUsecase) SendMessage(event MessageEvent) {
	chatID := event.Message.ChatID

	chatInfo, ok := w.onlineChats[chatID]
	if !ok {
		return
	}

	for userID := range chatInfo.users {
		if ch, ok := w.onlineUsers[userID]; ok {
			ch <- AnyEvent{
				TypeOfEvent: NewMessage,
				Event:       event,
			}
		}
	}
}
