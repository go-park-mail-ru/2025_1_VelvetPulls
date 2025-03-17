package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChats(t *testing.T) {
	t.Run("Cookies do not exist", func(t *testing.T) {
		body := bytes.NewReader(make([]byte, 0))
		r := httptest.NewRequest("GET", "/api/chats/", body)
		w := httptest.NewRecorder()

		Chats(w, r)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})
}
