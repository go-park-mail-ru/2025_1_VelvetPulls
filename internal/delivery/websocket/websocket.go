package http

import (
	"log"
	"net/http"

	usecase "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type websocketController struct {
	websocketUsecase usecase.IWebsocketUsecase
}

func NewWebsocketController(r *mux.Router, websocketUsecase usecase.IWebsocketUsecase) {
	controller := &websocketController{
		websocketUsecase: websocketUsecase,
	}
	r.Handle("/ws", http.HandlerFunc(controller.WebsocketConnection))
}

func (c *websocketController) WebsocketConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade to websocket: %v", err)
		http.Error(w, "could not open websocket connection", http.StatusBadRequest)
		return
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v", err)
			break
		}

		log.Printf("received message: %s", string(message))

		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("error writing message: %v", err)
			break
		}
	}
}
