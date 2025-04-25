package http

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/delivery/proto"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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
}

func NewWebsocketController(r *mux.Router, sessionClient authpb.SessionServiceClient, websocketUsecase usecase.IWebsocketUsecase) {
	controller := &WebsocketController{
		sessionClient:    sessionClient,
		websocketUsecase: websocketUsecase,
	}
	go websocketUsecase.ConsumeChats()
	go websocketUsecase.ConsumeMessages()

	r.HandleFunc("/ws", middleware.AuthMiddlewareWS(sessionClient)(controller.WebsocketConnection)).Methods("GET")
}

func (c *WebsocketController) WebsocketConnection(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Upgrade error:", zap.Error(err))
		return
	}

	userID := utils.GetUserIDFromCtx(r.Context())
	eventChan := make(chan model.AnyEvent, 100)

	// Регистрируем канал пользователя
	err = c.websocketUsecase.RegisterUserChannel(userID, eventChan)
	if err != nil {
		logger.Error("Broker Init error:", zap.Error(err))
		conn.Close()
		return
	}

	// Канал для отслеживания закрытия соединения
	done := make(chan struct{})

	// Горутина для чтения сообщений (чтобы обнаружить закрытие соединения)
	go func() {
		defer close(done)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	// Основной цикл отправки сообщений
	duration := 500 * time.Millisecond
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			// Соединение закрыто
			c.websocketUsecase.UnregisterUserChannel(userID, eventChan)
			conn.Close()
			return

		case message := <-eventChan:
			logger.Debug("Message delivery websocket: получены новые сообщения")

			if s, ok := message.Event.(interface{ Sanitize() }); ok {
				s.Sanitize()
			}

			if err := conn.WriteJSON(message.Event); err != nil {
				logger.Error("Write error:", zap.Error(err))
				c.websocketUsecase.UnregisterUserChannel(userID, eventChan)
				conn.Close()
				return
			}

		case <-ticker.C:
			// Периодическая проверка состояния
		}
	}
}
