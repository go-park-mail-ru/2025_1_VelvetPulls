package http

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/utils"
	"github.com/gorilla/mux"
)

type chatController struct {
	chatUsecase usecase.IChatUsecase
}

func NewChatController(r *mux.Router, chatUsecase usecase.IChatUsecase) {
	controller := &chatController{
		chatUsecase: chatUsecase,
	}

	r.HandleFunc("/chats/", controller.Chats).Methods(http.MethodGet)
}

// Chats возвращает чаты пользователя по сессии
// @Summary Получение чатов пользователя
// @Description Возвращает список чатов пользователя, ассоциированных с текущей сессией
// @Tags Chat
// @Accept json
// @Produce json
// @Param Cookie header string  false "token" default(token=xxx)
// @Success 200 {object} utils.JSONResponse
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/chats/ [get]
func (c *chatController) Chats(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetSessionCookie(r)
	if err != nil {
		utils.SendJSONResponse(w, http.StatusBadRequest, "invalid session token", false)
		return
	}

	chats, err := c.chatUsecase.FetchChatsBySession(token)
	if err != nil {
		if errors.Is(err, apperrors.ErrSessionNotFound) {
			utils.SendJSONResponse(w, http.StatusUnauthorized, "session not found", false)
			return
		}

		utils.SendJSONResponse(w, http.StatusInternalServerError, "internal server error", false)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, chats, true)
}
