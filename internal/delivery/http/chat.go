package http

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/gorilla/mux"
)

type chatController struct {
	chatUsecase    usecase.IChatUsecase
	sessionUsecase usecase.ISessionUsecase
}

func NewChatController(r *mux.Router, chatUsecase usecase.IChatUsecase, sessionUsecase usecase.ISessionUsecase) {
	controller := &chatController{
		chatUsecase:    chatUsecase,
		sessionUsecase: sessionUsecase,
	}

	r.Handle("/chats/", middleware.AuthMiddleware(sessionUsecase)(http.HandlerFunc(controller.Chats))).Methods(http.MethodGet)
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
	userID := r.Context().Value(utils.USER_ID_KEY).(string)

	chats, err := c.chatUsecase.FetchChatsByUserID(r.Context(), userID)
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
