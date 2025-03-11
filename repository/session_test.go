package repository

import (
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/stretchr/testify/require"
)

func TestSession(t *testing.T) {
	t.Run("Test sessions", func(t *testing.T) {
		sessionId := "abcdefg"
		var expectedSession *model.Session = nil

		// Сессии не существует

		session, err := GetSessionBySessId(sessionId)
		require.Equal(t, expectedSession, session)
		require.Equal(t, apperrors.ErrSessionNotFound, err)

		err = DeleteSession(sessionId)
		require.Equal(t, apperrors.ErrSessionNotFound, err)

		// Создаём сессию

		sessionId, err = CreateSession("ruslantus228")
		require.Equal(t, nil, err)

		// Сессия существует

		session, err = GetSessionBySessId(sessionId)
		require.Equal(t, sessions[sessionId], session)
		require.Equal(t, nil, err)

		err = DeleteSession(sessionId)
		require.Equal(t, nil, err)
	})
}
