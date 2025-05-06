package usecase_test

// // Тест получения сообщений чата (успешный сценарий) без участия WebSocket
// func TestGetChatMessages_NoWebsocket(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockMsgRepo := mocks.NewMockIMessageRepo(ctrl)
// 	mockChatRepo := mocks.NewMockIChatRepo(ctrl)
// 	// Передаем nil в качестве реализации IWebsocketUsecase
// 	msgUC := usecase.NewMessageUsecase(mockMsgRepo, mockChatRepo, nil)

// 	// Создаем контекст с логгером
// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, utils.LOGGER_ID_KEY, zap.NewNop())

// 	userID := uuid.New()
// 	chatID := uuid.New()

// 	// Ожидаем, что чат-репозиторий вернет допустимую роль для пользователя
// 	mockChatRepo.EXPECT().
// 		GetUserRoleInChat(ctx, userID, chatID).
// 		Return("owner", nil)

// 	// Подготавливаем список сообщений
// 	expectedMessages := []model.Message{
// 		{
// 			ID:     uuid.New(),
// 			ChatID: chatID,
// 			UserID: userID,
// 			Body:   "Hello",
// 			SentAt: time.Now(),
// 		},
// 		{
// 			ID:     uuid.New(),
// 			ChatID: chatID,
// 			UserID: userID,
// 			Body:   "World",
// 			SentAt: time.Now(),
// 		},
// 	}

// 	mockMsgRepo.EXPECT().
// 		GetMessages(ctx, chatID).
// 		Return(expectedMessages, nil)

// 	messages, err := msgUC.GetChatMessages(ctx, userID, chatID)
// 	require.NoError(t, err)
// 	assert.Equal(t, expectedMessages, messages)
// }

// // Тест успешной отправки сообщения без участия WebSocket (wsUsecase == nil)
// func TestSendMessage_NoWebsocket(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockMsgRepo := mocks.NewMockIMessageRepo(ctrl)
// 	mockChatRepo := mocks.NewMockIChatRepo(ctrl)
// 	// Передаем nil для wsUsecase – тогда метод SendMessage просто залогирует предупреждение, и событие не отправится.
// 	msgUC := usecase.NewMessageUsecase(mockMsgRepo, mockChatRepo, nil)

// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, utils.LOGGER_ID_KEY, zap.NewNop())

// 	userID := uuid.New()
// 	chatID := uuid.New()
// 	messageText := "Test message"
// 	messageInput := &model.MessageInput{
// 		Message: messageText,
// 	}

// 	// Возвращаем допустимую роль
// 	mockChatRepo.EXPECT().
// 		GetUserRoleInChat(ctx, userID, chatID).
// 		Return("member", nil)

// 	// Ожидаем создание сообщения
// 	createdMessage := &model.Message{
// 		ID:     uuid.New(),
// 		ChatID: chatID,
// 		UserID: userID,
// 		Body:   messageText,
// 		SentAt: time.Now(),
// 	}
// 	mockMsgRepo.EXPECT().
// 		CreateMessage(ctx, gomock.Any()).
// 		DoAndReturn(func(ctx context.Context, m *model.Message) (*model.Message, error) {
// 			// Можно проверить поля сообщения, если необходимо
// 			if m.ChatID != chatID || m.UserID != userID || m.Body != messageText {
// 				return nil, errors.New("invalid message input")
// 			}
// 			return createdMessage, nil
// 		})

// 	err := msgUC.SendMessage(ctx, messageInput, userID, chatID)
// 	assert.NoError(t, err)
// }

// // Тест отправки невалидного сообщения (например, пустое сообщение)
// func TestSendMessage_InvalidPayload(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockMsgRepo := mocks.NewMockIMessageRepo(ctrl)
// 	mockChatRepo := mocks.NewMockIChatRepo(ctrl)
// 	msgUC := usecase.NewMessageUsecase(mockMsgRepo, mockChatRepo, nil)

// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, utils.LOGGER_ID_KEY, zap.NewNop())

// 	userID := uuid.New()
// 	chatID := uuid.New()

// 	// Невалидный payload – пустое сообщение
// 	messageInput := &model.MessageInput{Message: ""}

// 	// Ожидаем, что роль верна, но валидация должна не пройти
// 	mockChatRepo.EXPECT().
// 		GetUserRoleInChat(ctx, userID, chatID).
// 		Return("owner", nil)

// 	err := msgUC.SendMessage(ctx, messageInput, userID, chatID)
// 	require.Error(t, err)
// 	assert.Contains(t, err.Error(), "invalid message input")
// }
