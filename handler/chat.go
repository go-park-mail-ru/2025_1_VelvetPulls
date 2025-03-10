package handler

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/service"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/utils"
)

// Chats возвращает чаты пользователя по сессии
// @Summary Получение чатов пользователя
// @Description Возвращает список чатов пользователя, ассоциированных с текущей сессией
// @Tags Chat
// @Accept json
// @Produce json
// @Param Cookie header string  false "token" default(token=xxx)
// @Success 200 {array} model.Chat
// @Failure 400 {object} utils.JSONResponse
// @Failure 500 {object} utils.JSONResponse
// @Router /api/chats/ [get]
func Chats(w http.ResponseWriter, r *http.Request) {
	token, err := utils.GetSessionCookie(r)
	if err != nil {
		// Отправляем ошибку с кодом 400 (неверный токен)
		_ = utils.SendJSONResponse(w, http.StatusBadRequest, err.Error(), false)
		return
	}

	// Получаем чаты по сессии
	chats, err := service.FetchChatsBySession(token)
	if err != nil {
		// Ошибка при получении чатов - отдаем ошибку и код
		_ = utils.SendJSONResponse(w, http.StatusInternalServerError, err.Error(), false)
		return
	}

	// Успешный ответ
	_ = utils.SendJSONResponse(w, http.StatusOK, chats, true)
}
