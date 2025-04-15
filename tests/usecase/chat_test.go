package usecase

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	mocks "github.com/go-park-mail-ru/2025_1_VelvetPulls/tests/mock"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestGetChatByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	chatID := uuid.New()
	logger := zap.NewNop() // безопасный заглушечный логгер
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)
	expectedChat := &model.Chat{
		ID:    chatID,
		Title: "Test Chat",
		// остальные поля можно опустить
	}

	// Указываем поведение мока
	mockRepo.EXPECT().
		GetChats(ctx, chatID).
		Return([]model.Chat{*expectedChat}, chatID, nil)

	usecase := usecase.NewChatUsecase(mockRepo)
	chat, err := usecase.GetChats(ctx, chatID)

	// Проверяем результат
	assert.NoError(t, err)
	assert.Equal(t, []model.Chat{*expectedChat}, chat)
}
