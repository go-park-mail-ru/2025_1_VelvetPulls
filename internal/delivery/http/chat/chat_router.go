package chat

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase/chat"
	"github.com/gorilla/mux"
)

type chatController struct {
	chatUsecase chat.ChatUsecaseInterface
}

func NewChatController(r *mux.Router, chatUsecase chat.ChatUsecaseInterface) {
	controller := &chatController{
		chatUsecase: chatUsecase,
	}

	r.HandleFunc("/chats/", controller.Chats).Methods(http.MethodGet)
}
