package usecase

import (
	"encoding/json"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
)

const (
	NewMessage    = "newMessage"
	UpdateMessage = "updateMessage"
	DeleteMessage = "deleteMessage"
)

func SerializeMessageEvent(event model.MessageEvent) ([]byte, error) {
	return json.Marshal(event)
}

func DeserializeMessageEvent(data []byte) (model.MessageEvent, error) {
	var event model.MessageEvent
	err := json.Unmarshal(data, &event)
	if err != nil {
		return model.MessageEvent{}, err
	}
	return event, nil
}

func (w *WebsocketUsecase) ConsumeMessages() {
	for msg := range w.messageChan {
		w.SendMessage(msg)
	}
}

func (w *WebsocketUsecase) SendMessage(event model.MessageEvent) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	chatID := event.Message.ChatID
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
