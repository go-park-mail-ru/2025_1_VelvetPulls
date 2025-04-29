package usecase

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/websocket_service/model"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// IWebsocketUsecase defines methods to manage WebSocket channels and NATS integration
type IWebsocketUsecase interface {
	RegisterUserChannel(userID uuid.UUID, eventChan chan model.AnyEvent) error
	UnregisterUserChannel(userID uuid.UUID, eventChan chan model.AnyEvent)
	GetUserChannels(userID uuid.UUID) []chan model.AnyEvent
	InitChat(chatID uuid.UUID, users []uuid.UUID) error
	PublishMessage(event model.MessageEvent) error
	PublishChatEvent(event model.ChatEvent) error
}

// WebsocketUsecase implements IWebsocketUsecase using NATS for transport
type WebsocketUsecase struct {
	nc            *nats.Conn
	onlineChats   map[uuid.UUID]model.ChatInfoWS      // chatID -> ChatInfoWS{Events, Users}
	onlineUsers   map[uuid.UUID][]chan model.AnyEvent // userID -> slice of event channels
	subscriptions map[uuid.UUID][]*nats.Subscription  // chatID -> NATS subscriptions
	mu            sync.RWMutex
}

// NewWebsocketUsecase creates a new WebsocketUsecase
func NewWebsocketUsecase(nc *nats.Conn) IWebsocketUsecase {
	return &WebsocketUsecase{
		nc:            nc,
		onlineChats:   make(map[uuid.UUID]model.ChatInfoWS),
		onlineUsers:   make(map[uuid.UUID][]chan model.AnyEvent),
		subscriptions: make(map[uuid.UUID][]*nats.Subscription),
	}
}

// RegisterUserChannel adds a user's event channel
func (w *WebsocketUsecase) RegisterUserChannel(userID uuid.UUID, eventChan chan model.AnyEvent) error {
	w.mu.Lock()
	w.onlineUsers[userID] = append(w.onlineUsers[userID], eventChan)
	w.mu.Unlock()
	return nil
}

// UnregisterUserChannel removes a user's event channel
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

// GetUserChannels returns a copy of channels for a user
func (w *WebsocketUsecase) GetUserChannels(userID uuid.UUID) []chan model.AnyEvent {
	w.mu.RLock()
	orig := w.onlineUsers[userID]
	w.mu.RUnlock()
	copy := make([]chan model.AnyEvent, len(orig))
	copy = append(copy[:0], orig...)
	return copy
}

// InitChat subscribes to NATS subjects for a chat and tracks membership
func (w *WebsocketUsecase) InitChat(chatID uuid.UUID, users []uuid.UUID) error {
	w.mu.Lock()
	// Initialize chat info if not exists
	if _, exists := w.onlineChats[chatID]; !exists {
		eventsChan := make(chan model.Event, 100)
		w.onlineChats[chatID] = model.ChatInfoWS{Events: eventsChan, Users: make(map[uuid.UUID]struct{})}

		// Start goroutine to handle membership events
		go w.handleMembershipEvents(chatID)

		// Subscribe to membership events
		subM, err := w.nc.Subscribe(fmt.Sprintf("chat.%s.events", chatID.String()), func(msg *nats.Msg) {
			var ev model.Event
			if err := json.Unmarshal(msg.Data, &ev); err != nil {
				zap.L().Error("failed to unmarshal membership event", zap.Error(err))
				return
			}
			w.onlineChats[chatID].Events <- ev
		})
		if err != nil {
			w.mu.Unlock()
			return fmt.Errorf("subscribe membership: %w", err)
		}

		// Subscribe to chat messages
		subMsg, err := w.nc.Subscribe(fmt.Sprintf("chat.%s.messages", chatID.String()), func(msg *nats.Msg) {
			var me model.MessageEvent
			if err := json.Unmarshal(msg.Data, &me); err != nil {
				zap.L().Error("failed to unmarshal message event", zap.Error(err))
				return
			}
			w.SendMessage(me)
		})
		if err != nil {
			w.mu.Unlock()
			return fmt.Errorf("subscribe messages: %w", err)
		}

		// Save subscriptions
		w.subscriptions[chatID] = []*nats.Subscription{subM, subMsg}
	}
	// Add initial users
	for _, uid := range users {
		w.onlineChats[chatID].Users[uid] = struct{}{}
	}
	w.mu.Unlock()
	return nil
}

// PublishMessage publishes a message event to NATS
func (w *WebsocketUsecase) PublishMessage(event model.MessageEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	subj := fmt.Sprintf("chat.%s.messages", event.Message.ChatID.String())
	return w.nc.Publish(subj, data)
}

// PublishChatEvent publishes a membership event to NATS
func (w *WebsocketUsecase) PublishChatEvent(event model.ChatEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	subj := fmt.Sprintf("chat.%s.events", event.Chat.ID.String())
	return w.nc.Publish(subj, data)
}

// handleMembershipEvents processes Event from NATS and updates onlineChats
func (w *WebsocketUsecase) handleMembershipEvents(chatID uuid.UUID) {
	ch := w.onlineChats[chatID].Events
	for ev := range ch {
		w.mu.Lock()
		chatInfo := w.onlineChats[chatID]
		switch ev.Action {
		case utils.AddWebsocketUser:
			for _, uid := range ev.Users {
				chatInfo.Users[uid] = struct{}{}
			}
		case utils.DeleteWebsocketUser:
			for _, uid := range ev.Users {
				delete(chatInfo.Users, uid)
			}
		}
		w.onlineChats[chatID] = chatInfo
		w.mu.Unlock()
	}
}

// SendMessage distributes MessageEvent to user channels
func (w *WebsocketUsecase) SendMessage(event model.MessageEvent) {
	w.mu.RLock()
	chatInfo, ok := w.onlineChats[event.Message.ChatID]
	if !ok {
		w.mu.RUnlock()
		return
	}
	users := make([]uuid.UUID, 0, len(chatInfo.Users))
	for uid := range chatInfo.Users {
		users = append(users, uid)
	}
	chans := w.onlineUsers
	w.mu.RUnlock()

	for _, userID := range users {
		for _, ch := range chans[userID] {
			select {
			case ch <- model.AnyEvent{TypeOfEvent: event.Action, Event: event}:
			default:
			}
		}
	}
}
