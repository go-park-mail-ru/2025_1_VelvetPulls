package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	mocks "github.com/go-park-mail-ru/2025_1_VelvetPulls/tests/usecase/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

// Тест получения сообщений чата (успешный сценарий) без участия WebSocket
func TestGetChatMessages_NoWebsocket(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMsgRepo := mocks.NewMockIMessageRepo(ctrl)
	mockChatRepo := mocks.NewMockIChatRepo(ctrl)
	// Передаем nil в качестве реализации IWebsocketUsecase
	msgUC := usecase.NewMessageUsecase(mockMsgRepo, mockChatRepo, nil)

	// Создаем контекст с логгером
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.LOGGER_ID_KEY, zap.NewNop())

	userID := uuid.New()
	chatID := uuid.New()

	// Ожидаем, что чат-репозиторий вернет допустимую роль для пользователя
	mockChatRepo.EXPECT().
		GetUserRoleInChat(ctx, userID, chatID).
		Return("owner", nil)

	// Подготавливаем список сообщений
	expectedMessages := []model.Message{
		{
			ID:     uuid.New(),
			ChatID: chatID,
			UserID: userID,
			Body:   "Hello",
			SentAt: time.Now(),
		},
		{
			ID:     uuid.New(),
			ChatID: chatID,
			UserID: userID,
			Body:   "World",
			SentAt: time.Now(),
		},
	}

	mockMsgRepo.EXPECT().
		GetMessages(ctx, chatID).
		Return(expectedMessages, nil)

	messages, err := msgUC.GetChatMessages(ctx, userID, chatID)
	require.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
}
