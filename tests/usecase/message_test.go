package usecase_test

// func TestSendMessage_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockMsgRepo := mocks.NewMockIMessageRepo(ctrl)
// 	mockChatRepo := mocks.NewMockIChatRepo(ctrl)

// 	msgUC := usecase.NewMessageUsecase(mockMsgRepo, mockChatRepo, nil)

// 	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())
// 	userID := uuid.New()
// 	chatID := uuid.New()

// 	input := &model.MessageInput{Message: "test message"}

// 	mockChatRepo.EXPECT().
// 		GetChatByID(ctx, chatID).
// 		Return(&model.Chat{ID: chatID, Type: string(model.ChatTypeGroup)}, nil)

// 	mockChatRepo.EXPECT().
// 		GetUserRoleInChat(ctx, userID, chatID).
// 		Return("member", nil)

// 	savedMsg := &model.Message{
// 		ID:     uuid.New(),
// 		ChatID: chatID,
// 		UserID: userID,
// 		Body:   input.Message,
// 		SentAt: time.Now(),
// 	}
// 	mockMsgRepo.EXPECT().
// 		CreateMessage(ctx, gomock.Any()).
// 		Return(savedMsg, nil)

// 	err := msgUC.SendMessage(ctx, input, userID, chatID)
// 	require.NoError(t, err)
// }

// func TestUpdateMessage_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockMsgRepo := mocks.NewMockIMessageRepo(ctrl)
// 	mockChatRepo := mocks.NewMockIChatRepo(ctrl)

// 	msgUC := usecase.NewMessageUsecase(mockMsgRepo, mockChatRepo, nil)

// 	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())
// 	userID := uuid.New()
// 	chatID := uuid.New()
// 	messageID := uuid.New()
// 	input := &model.MessageInput{Message: "updated"}

// 	mockChatRepo.EXPECT().
// 		GetUserRoleInChat(ctx, userID, chatID).
// 		Return("member", nil)

// 	original := &model.Message{
// 		ID:     messageID,
// 		ChatID: chatID,
// 		UserID: userID,
// 		Body:   "old",
// 		SentAt: time.Now(),
// 	}

// 	mockMsgRepo.EXPECT().
// 		GetMessage(ctx, messageID).
// 		Return(original, nil)

// 	updated := &model.Message{
// 		ID:     messageID,
// 		ChatID: chatID,
// 		UserID: userID,
// 		Body:   "updated",
// 		SentAt: time.Now(),
// 	}

// 	mockMsgRepo.EXPECT().
// 		UpdateMessage(ctx, messageID, input.Message).
// 		Return(updated, nil)

// 	err := msgUC.UpdateMessage(ctx, messageID, input, userID, chatID)
// 	require.NoError(t, err)
// }

// func TestDeleteMessage_Success(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockMsgRepo := mocks.NewMockIMessageRepo(ctrl)
// 	mockChatRepo := mocks.NewMockIChatRepo(ctrl)

// 	msgUC := usecase.NewMessageUsecase(mockMsgRepo, mockChatRepo, nil)

// 	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())
// 	userID := uuid.New()
// 	chatID := uuid.New()
// 	messageID := uuid.New()

// 	mockChatRepo.EXPECT().
// 		GetUserRoleInChat(ctx, userID, chatID).
// 		Return("member", nil)

// 	msg := &model.Message{
// 		ID:     messageID,
// 		ChatID: chatID,
// 		UserID: userID,
// 		Body:   "message to delete",
// 		SentAt: time.Now(),
// 	}
// 	mockMsgRepo.EXPECT().
// 		GetMessage(ctx, messageID).
// 		Return(msg, nil)

// 	mockMsgRepo.EXPECT().
// 		DeleteMessage(ctx, messageID).
// 		Return(msg, nil)

// 	err := msgUC.DeleteMessage(ctx, messageID, userID, chatID)
// 	require.NoError(t, err)
// }
