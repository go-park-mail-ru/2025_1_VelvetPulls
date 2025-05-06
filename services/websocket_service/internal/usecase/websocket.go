package usecase

import (
	"sync"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/websocket_service/internal/model"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type IWebsocketUsecase interface {
	RegisterUserChannel(userID uuid.UUID, eventChan chan model.AnyEvent) error
	UnregisterUserChannel(userID uuid.UUID, eventChan chan model.AnyEvent)
	GetUserChannels(userID uuid.UUID) []chan model.AnyEvent
}

type WebsocketUsecase struct {
	nc            *nats.Conn
	onlineChats   map[uuid.UUID]model.ChatInfoWS      // chatID -> ChatInfoWS{Events, Users}
	onlineUsers   map[uuid.UUID][]chan model.AnyEvent // userID -> slice of event channels
	subscriptions map[uuid.UUID][]*nats.Subscription  // chatID -> NATS subscriptions
	mu            sync.RWMutex
}

func NewWebsocketUsecase(nc *nats.Conn) IWebsocketUsecase {
	return &WebsocketUsecase{
		nc:            nc,
		onlineChats:   make(map[uuid.UUID]model.ChatInfoWS),
		onlineUsers:   make(map[uuid.UUID][]chan model.AnyEvent),
		subscriptions: make(map[uuid.UUID][]*nats.Subscription),
	}
}

func (w *WebsocketUsecase) RegisterUserChannel(userID uuid.UUID, eventChan chan model.AnyEvent) error {
	w.mu.Lock()
	w.onlineUsers[userID] = append(w.onlineUsers[userID], eventChan)
	w.mu.Unlock()
	return nil
}

func (w *WebsocketUsecase) UnregisterUserChannel(userID uuid.UUID, eventChan chan model.AnyEvent) {
	w.mu.Lock()
	channels := w.onlineUsers[userID]
	for i, ch := range channels {
		if ch == eventChan {
			w.onlineUsers[userID] = append(channels[:i], channels[i+1:]...)
			break
		}
	}
	if len(w.onlineUsers[userID]) == 0 {
		delete(w.onlineUsers, userID)
	}
	w.mu.Unlock()
	close(eventChan)
}

func (w *WebsocketUsecase) GetUserChannels(userID uuid.UUID) []chan model.AnyEvent {
	w.mu.RLock()
	orig := w.onlineUsers[userID]
	w.mu.RUnlock()
	copy := make([]chan model.AnyEvent, len(orig))
	copy = append(copy[:0], orig...)
	return copy
}
