package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/proto"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/websocket_service/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/websocket_service/internal/usecase"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketController struct {
	sessionClient    authpb.SessionServiceClient
	websocketUsecase usecase.IWebsocketUsecase
	nc               *nats.Conn
}

func NewWebsocketController(r *mux.Router, sessionClient authpb.SessionServiceClient, websocketUsecase usecase.IWebsocketUsecase, nc *nats.Conn) {
	controller := &WebsocketController{
		sessionClient:    sessionClient,
		websocketUsecase: websocketUsecase,
		nc:               nc,
	}

	r.HandleFunc("/ws", middleware.AuthMiddlewareWS(sessionClient)(controller.WebsocketConnection)).Methods("GET")
}
func (c *WebsocketController) WebsocketConnection(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("WebSocket upgrade failed", zap.Error(err))
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())
	eventChan := make(chan model.AnyEvent, 100)

	if err := c.websocketUsecase.RegisterUserChannel(userID, eventChan); err != nil {
		logger.Error("RegisterUserChannel failed", zap.Error(err))
		conn.Close()
		return
	}
	defer c.websocketUsecase.UnregisterUserChannel(userID, eventChan)

	// Подписка на персональные события
	subUser, err := c.nc.Subscribe(fmt.Sprintf("user.%s.events", userID.String()), func(msg *nats.Msg) {
		var anyEv model.AnyEvent
		if err := json.Unmarshal(msg.Data, &anyEv); err != nil {
			logger.Error("unmarshal user event", zap.Error(err))
			return
		}
		// Обработка события для конкретного пользователя
		c.handleUserEvent(anyEv, eventChan)
	})
	if err != nil {
		logger.Error("Subscribe user.* failed", zap.Error(err))
		conn.Close()
		return
	}
	defer func() {
		if err := subUser.Unsubscribe(); err != nil {
			logger.Error("Failed to unsubscribe user subscription", zap.Error(err))
		}
	}()

	// Подписка на все события чата через wildcard
	subChat, err := c.nc.Subscribe("chat.*.*", func(msg *nats.Msg) {
		parts := strings.Split(msg.Subject, ".")
		if len(parts) != 3 {
			return
		}
		kind := parts[2] // "messages" или "events"
		switch kind {
		case "messages":
			var me model.MessageEvent
			if err := json.Unmarshal(msg.Data, &me); err != nil {
				logger.Error("unmarshal message event", zap.Error(err))
				return
			}
			c.handleMessageEvent(me, eventChan)

		case "events":
			var ce model.ChatEvent
			if err := json.Unmarshal(msg.Data, &ce); err != nil {
				logger.Error("unmarshal chat event", zap.Error(err))
				return
			}
			c.handleChatEvent(ce, eventChan)
		}
	})
	if err != nil {
		logger.Error("Subscribe chat.*.* failed", zap.Error(err))
		conn.Close()
		return
	}
	defer func() {
		if err := subChat.Unsubscribe(); err != nil {
			logger.Error("Failed to unsubscribe chat subscription", zap.Error(err))
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			conn.Close()
			return

		case anyEv := <-eventChan:
			if s, ok := anyEv.Event.(interface{ Sanitize() }); ok {
				s.Sanitize()
			}
			if err := conn.WriteJSON(anyEv.Event); err != nil {
				logger.Error("WriteJSON failed", zap.Error(err))
				conn.Close()
				return
			}

		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				conn.Close()
				return
			}
		}
	}
}

func (c *WebsocketController) handleUserEvent(event model.AnyEvent, eventChan chan model.AnyEvent) {
	eventChan <- event
}

func (c *WebsocketController) handleMessageEvent(event model.MessageEvent, eventChan chan model.AnyEvent) {
	eventChan <- model.AnyEvent{TypeOfEvent: event.Action, Event: event}
}

func (c *WebsocketController) handleChatEvent(event model.ChatEvent, eventChan chan model.AnyEvent) {
	eventChan <- model.AnyEvent{TypeOfEvent: event.Action, Event: event}
}
