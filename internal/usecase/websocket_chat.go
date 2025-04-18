package usecase

import (
	"encoding/json"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
)

const (
	NewChat     = "newChat"
	UpdateChat  = "updateChat"
	DeleteChat  = "deleteChat"
	AddUsers    = "addUsers"
	RemoveUsers = "removeUsers"
)

func SerializeChatEvent(event model.ChatEvent) ([]byte, error) {
	return json.Marshal(event)
}

func DeserializeChatEvent(data []byte) (model.ChatEvent, error) {
	var event model.ChatEvent
	err := json.Unmarshal(data, &event)
	if err != nil {
		return model.ChatEvent{}, err
	}
	return event, nil
}

func (w *WebsocketUsecase) ConsumeChats() {
	for event := range w.chatChan {
		w.SendChatEvent(event)
	}
}

func (w *WebsocketUsecase) SendChatEvent(event model.ChatEvent) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	chatID := event.Chat.ID
	chatInfo, ok := w.onlineChats[chatID]
	if !ok {
		return
	}
	for userID := range chatInfo.Users {
		if chans, ok := w.onlineUsers[userID]; ok {
			for _, ch := range chans {
				select {
				case ch <- model.AnyEvent{
					TypeOfEvent: event.Action,
					Event:       event,
				}:
				default:
					// Если канал полон, пропускаем сообщение
				}
			}
		}
	}
}
