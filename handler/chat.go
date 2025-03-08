package handler

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/service"
	utils "github.com/go-park-mail-ru/2025_1_VelvetPulls/utils"
)

// Chats возвращает чаты пользователя по сессии
// @Summary Получение чатов пользователя
// @Description Возвращает список чатов пользователя, ассоциированных с текущей сессией
// @Tags Chat
// @Accept json
// @Produce json
// @Param Cookie header string  false "token" default(token=xxx)
// @Success 200 {array} model.Chat
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /api/chats/ [get]
func Chats(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetSessionCookie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userResponse, err := service.FetchChatsBySession(token)
	if err != nil {
		http.Error(w, userResponse.Body.(error).Error(), userResponse.StatusCode)
		return
	}

	if err := utils.SendJSONResponse(w, userResponse.StatusCode, userResponse.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
