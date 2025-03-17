package service

import (
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/repository"
	"github.com/stretchr/testify/require"
)

func TestFetchChatsBySession(t *testing.T) {
	negativeTests := []struct {
		testName string
		token    string
		chats    []*model.Chat
		err      error
	}{
		{
			testName: "Session is not found",
			token:    "abcdefg",
			chats:    nil,
			err:      apperrors.ErrSessionNotFound,
		},
	}

	for _, tt := range negativeTests {
		t.Run(tt.testName, func(t *testing.T) {
			chats, err := FetchChatsBySession(tt.token)

			require.Equal(t, tt.chats, chats)
			require.Equal(t, tt.err, err)
		})
	}

	t.Run("Get chats", func(t *testing.T) {
		userName := "ruslantus228"
		token, err := repository.CreateSession(userName)
		require.Equal(t, nil, err)
		require.NotEqual(t, "", token)

		expectedChats, err := repository.GetChatsByUsername(userName)
		require.Equal(t, nil, err)

		actualChats, err := FetchChatsBySession(token)
		require.Equal(t, nil, err)
		require.Equal(t, expectedChats, actualChats)
	})
}
