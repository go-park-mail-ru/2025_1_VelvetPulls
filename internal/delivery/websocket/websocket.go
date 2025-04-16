package http

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
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
	sessionUsecase   usecase.ISessionUsecase
	websocketUsecase usecase.IWebsocketUsecase
}

func NewWebsocketController(r *mux.Router, sessionUsecase usecase.ISessionUsecase, websocketUsecase usecase.IWebsocketUsecase) {
	controller := &WebsocketController{
		sessionUsecase:   sessionUsecase,
		websocketUsecase: websocketUsecase,
	}

	go websocketUsecase.ConsumeMessages()

	r.HandleFunc("/ws", middleware.AuthMiddlewareWS(sessionUsecase)(controller.WebsocketConnection)).Methods("GET")
}

func (c *WebsocketController) WebsocketConnection(w http.ResponseWriter, r *http.Request) {
	logger := utils.GetLoggerFromCtx(r.Context())
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Upgrade error:", zap.Error(err))
		return
	}

	defer conn.Close()

	userID := utils.GetUserIDFromCtx(r.Context())
	eventChan := make(chan usecase.AnyEvent, 100)

	err = c.websocketUsecase.InitBrokersForUser(userID, eventChan)
	if err != nil {
		logger.Error("Broker Init error:", zap.Error(err))
		return
	}

	// пока соеденено
	duration := 500 * time.Millisecond

	for {
		select {
		case message := <-eventChan:
			logger.Info("Message delivery websocket: получены новые сообщения")

			if s, ok := message.Event.(interface{ Sanitize() }); ok {
				s.Sanitize()
			}

			conn.WriteJSON(message.Event)

		default:
			time.Sleep(duration)
		}
	}

}
