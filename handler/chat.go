package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/repository"
)

func Chats(w http.ResponseWriter, r *http.Request) {
	sessionId, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if _, err := repository.GetSessionByID(sessionId.Value); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userId, err := strconv.Atoi(sessionId.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chats := repository.GetChatsWithUser(int64(userId))

	result, err := json.Marshal(chats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, result)
}
