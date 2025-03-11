package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/service"
	"github.com/stretchr/testify/require"
)

func TestChats(t *testing.T) {
	// tests := []struct {
	// 	name   string
	// 	status int
	// }{
	// 	{
	// 		name:   "First test",
	// 		status: http.StatusBadRequest,
	// 	},
	// }

	// for _, test := range tests {
	// 	t.Run(test.name, func(t *testing.T) {
	// 		body := bytes.NewReader(make([]byte, 0))

	// 		r := httptest.NewRequest("GET", "/api/chats/", body)
	// 		w := httptest.NewRecorder()

	// 		Chats(w, r)

	// 		require.Equal(t, test.status, w.Code)
	// 	})
	// }

	t.Run("Cookies do not exist", func(t *testing.T) {
		body := bytes.NewReader(make([]byte, 0))
		r := httptest.NewRequest("GET", "/api/chats/", body)
		w := httptest.NewRecorder()

		Chats(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Successful GET", func(t *testing.T) {
		sessionId, err := service.LoginUser(model.LoginCredentials{
			Username: "ruslantus228",
			Password: "qwerty",
		})
		require.Equal(t, nil, err)
		require.NotEqual(t, "", sessionId)

		body := bytes.NewReader(make([]byte, 0))
		r := httptest.NewRequest("GET", "/api/chats/", body)
		// TODO: Проставить Cookie в r

		w := httptest.NewRecorder()

		Chats(w, r)

		require.Equal(t, http.StatusOK, w.Code)
	})
}
