package repository

import (
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/stretchr/testify/require"
)

func TestGetChatsByUsername(t *testing.T) {
	tests := []struct {
		testName      string
		userName      string
		expectedChats []*model.Chat
		expectedError error
	}{
		{
			testName:      "Get chats: ruslantus228",
			userName:      "ruslantus228",
			expectedChats: chats[:],
			expectedError: nil,
		},
		{
			testName:      "Get chats: ilyaaaaaaaaz",
			userName:      "ilyaaaaaaaaz",
			expectedChats: make([]*model.Chat, 0),
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			userChats, err := GetChatsByUsername(tt.userName)

			require.Equal(t, tt.expectedChats, userChats)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestAddChat(t *testing.T) {
	tests := []struct {
		testName string
		chat     *model.Chat
	}{
		{
			testName: "New chat",
			chat:     &model.Chat{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			AddChat(tt.chat)
		})
	}
}
