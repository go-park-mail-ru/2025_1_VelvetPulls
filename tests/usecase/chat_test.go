package usecase_test

import (
	"context"
	"testing"

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

func TestGetChatByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	chatID := uuid.New()
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)
	expectedChat := &model.Chat{
		ID:    chatID,
		Title: "Test Chat",
	}

	mockRepo.EXPECT().
		GetChats(ctx, chatID).
		Return([]model.Chat{*expectedChat}, chatID, nil)

	usecase := usecase.NewChatUsecase(mockRepo, nil)
	chat, err := usecase.GetChats(ctx, chatID)

	assert.NoError(t, err)
	assert.Equal(t, []model.Chat{*expectedChat}, chat)
}

func TestGetChatInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatRepo := mocks.NewMockIChatRepo(ctrl)
	chatUC := usecase.NewChatUsecase(mockChatRepo, nil)

	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	chatID := uuid.New()

	avatarPath := "chat_avatar.png"
	userAvatar := "user_avatar.jpg"
	chatType := "group"
	title := "Go Dev Chat"
	role := "owner"

	name := "Test User"

	userInChat := model.UserInChat{
		ID:         uuid.New(),
		Username:   "testuser",
		Name:       &name,
		AvatarPath: &userAvatar,
		Role:       &role,
	}

	expectedChat := &model.Chat{
		ID:         chatID,
		AvatarPath: &avatarPath,
		Type:       chatType,
		Title:      title,
	}

	mockChatRepo.EXPECT().
		GetUserRoleInChat(ctx, userID, chatID).
		Return(role, nil)

	mockChatRepo.EXPECT().
		GetChatByID(ctx, chatID).
		Return(expectedChat, nil)

	mockChatRepo.EXPECT().
		GetUsersFromChat(ctx, chatID).
		Return([]model.UserInChat{userInChat}, nil)

	result, err := chatUC.GetChatInfo(ctx, userID, chatID)

	require.NoError(t, err, "expected no error from GetChatInfo")
	require.NotNil(t, result, "expected non-nil result")

	assert.Equal(t, expectedChat.ID, result.ID)
	assert.Equal(t, expectedChat.AvatarPath, result.AvatarPath)
	assert.Equal(t, expectedChat.Type, result.Type)
	assert.Equal(t, expectedChat.Title, result.Title)
	require.Len(t, result.Users, 1)

	gotUser := result.Users[0]
	assert.Equal(t, userInChat.Username, gotUser.Username)
	assert.Equal(t, *userInChat.Name, *gotUser.Name)
	assert.Equal(t, userInChat.AvatarPath, gotUser.AvatarPath)
	assert.Equal(t, userInChat.Role, gotUser.Role)
}

func TestCreateGroupChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	chatID := uuid.New()
	createChat := &model.CreateChat{
		Title: "Group Chat",
		Type:  "group",
	}

	mockRepo.EXPECT().
		CreateChat(ctx, createChat).
		Return(chatID, "", nil)

	mockRepo.EXPECT().
		AddUserToChatByID(ctx, userID, "owner", chatID).
		Return(nil)

	mockRepo.EXPECT().
		GetUserRoleInChat(ctx, userID, chatID).
		Return("owner", nil)

	mockRepo.EXPECT().
		GetChatByID(ctx, chatID).
		Return(&model.Chat{ID: chatID, Type: "group", Title: "Group Chat"}, nil)

	mockRepo.EXPECT().
		GetUsersFromChat(ctx, chatID).
		Return([]model.UserInChat{{ID: userID, Username: "user1"}}, nil)

	uc := usecase.NewChatUsecase(mockRepo, nil)
	info, err := uc.CreateChat(ctx, userID, createChat)

	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, chatID, info.ID)
	assert.Equal(t, "Group Chat", info.Title)
}

func TestUpdateChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	chatID := uuid.New()
	title := "Updated Title"
	updateChat := &model.UpdateChat{
		ID:    chatID,
		Title: &title,
	}

	mockRepo.EXPECT().GetChatByID(gomock.Any(), chatID).
		Return(&model.Chat{ID: chatID, Type: "group"}, nil).
		Times(1)

	mockRepo.EXPECT().GetUserRoleInChat(gomock.Any(), userID, chatID).
		Return("owner", nil).
		Times(1)

	mockRepo.EXPECT().UpdateChat(gomock.Any(), updateChat).
		Return("", "", nil).
		Times(1)

	mockRepo.EXPECT().GetUserRoleInChat(gomock.Any(), userID, chatID).
		Return("owner", nil).
		Times(1)

	mockRepo.EXPECT().GetChatByID(gomock.Any(), chatID).
		Return(&model.Chat{ID: chatID, Type: "group", Title: title}, nil).
		Times(1)

	mockRepo.EXPECT().GetUsersFromChat(gomock.Any(), chatID).
		Return([]model.UserInChat{}, nil).
		Times(1)

	uc := usecase.NewChatUsecase(mockRepo, nil)
	info, err := uc.UpdateChat(ctx, userID, updateChat)

	assert.NoError(t, err)
	assert.Equal(t, title, info.Title)
}

func TestDeleteChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	chatID := uuid.New()

	mockRepo.EXPECT().GetUserRoleInChat(ctx, userID, chatID).Return("owner", nil)
	mockRepo.EXPECT().DeleteChat(ctx, chatID).Return(nil)

	usecase := usecase.NewChatUsecase(mockRepo, nil)
	err := usecase.DeleteChat(ctx, userID, chatID)

	assert.NoError(t, err)
}

func TestAddUsersIntoChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	chatID := uuid.New()
	usernames := []string{"alice", "bob"}

	mockRepo.EXPECT().GetChatByID(ctx, chatID).Return(&model.Chat{ID: chatID, Type: "group"}, nil)
	mockRepo.EXPECT().GetUserRoleInChat(ctx, userID, chatID).Return("owner", nil)
	mockRepo.EXPECT().AddUserToChatByUsername(ctx, "alice", "member", chatID).Return(nil)
	mockRepo.EXPECT().AddUserToChatByUsername(ctx, "bob", "member", chatID).Return(nil)

	usecase := usecase.NewChatUsecase(mockRepo, nil)
	res, err := usecase.AddUsersIntoChat(ctx, userID, usernames, chatID)

	assert.NoError(t, err)
	assert.ElementsMatch(t, usernames, res.AddedUsers)
	assert.Empty(t, res.NotAddedUsers)
}

func TestDeleteUserFromChat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	chatID := uuid.New()
	toDelete := []string{"charlie"}

	mockRepo.EXPECT().GetChatByID(ctx, chatID).Return(&model.Chat{ID: chatID, Type: "group"}, nil)
	mockRepo.EXPECT().GetUserRoleInChat(ctx, userID, chatID).Return("owner", nil)
	mockRepo.EXPECT().RemoveUserFromChatByUsername(ctx, "charlie", chatID).Return(nil)

	usecase := usecase.NewChatUsecase(mockRepo, nil)
	res, err := usecase.DeleteUserFromChat(ctx, userID, toDelete, chatID)

	assert.NoError(t, err)
	assert.Equal(t, toDelete, res.DeletedUsers)
}

func TestGetChats_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)
	userID := uuid.New()

	mockRepo.EXPECT().
		GetChats(ctx, userID).
		Return(nil, uuid.Nil, assert.AnError)

	uc := usecase.NewChatUsecase(mockRepo, nil)
	_, err := uc.GetChats(ctx, userID)

	assert.Error(t, err)
}

func TestGetChatInfo_PermissionDenied(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	chatID := uuid.New()

	mockRepo.EXPECT().
		GetUserRoleInChat(ctx, userID, chatID).
		Return("", nil) // Пустая роль - нет доступа

	uc := usecase.NewChatUsecase(mockRepo, nil)
	_, err := uc.GetChatInfo(ctx, userID, chatID)

	assert.Error(t, err)
	assert.Equal(t, usecase.ErrPermissionDenied, err)
}

func TestCreateDialogChat_ExistingDialog(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	dialogUser := "existing_user"
	existingChatID := uuid.New()

	createChat := &model.CreateChat{
		Type:       "dialog",
		Title:      "Dialog",
		DialogUser: dialogUser,
	}

	// Мокаем получение чатов пользователя
	existingChat := model.Chat{
		ID:   existingChatID,
		Type: "dialog",
	}
	mockRepo.EXPECT().
		GetChats(ctx, userID).
		Return([]model.Chat{existingChat}, userID, nil)

	// Мокаем получение пользователей из существующего чата
	mockRepo.EXPECT().
		GetUsersFromChat(ctx, existingChatID).
		Return([]model.UserInChat{
			{ID: userID, Username: "current_user"},
			{Username: dialogUser},
		}, nil)

	// Ожидаем получение информации о существующем чате
	mockRepo.EXPECT().
		GetUserRoleInChat(ctx, userID, existingChatID).
		Return("owner", nil)
	mockRepo.EXPECT().
		GetChatByID(ctx, existingChatID).
		Return(&model.Chat{ID: existingChatID, Type: "dialog"}, nil)
	mockRepo.EXPECT().
		GetUsersFromChat(ctx, existingChatID).
		Return([]model.UserInChat{}, nil)

	uc := usecase.NewChatUsecase(mockRepo, nil)
	info, err := uc.CreateChat(ctx, userID, createChat)

	assert.NoError(t, err)
	assert.Equal(t, existingChatID, info.ID)
}

func TestCreateChat_ValidationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	invalidChat := &model.CreateChat{
		Type:  "invalid_type", // Невалидный тип
		Title: "Chat",
	}

	uc := usecase.NewChatUsecase(mockRepo, nil)
	_, err := uc.CreateChat(ctx, userID, invalidChat)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation")
}

func TestUpdateChat_NotOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	chatID := uuid.New()
	title := "New Title"
	updateChat := &model.UpdateChat{
		ID:    chatID,
		Title: &title,
	}

	mockRepo.EXPECT().GetChatByID(ctx, chatID).
		Return(&model.Chat{ID: chatID, Type: "group"}, nil)
	mockRepo.EXPECT().GetUserRoleInChat(ctx, userID, chatID).
		Return("member", nil) // Не owner

	uc := usecase.NewChatUsecase(mockRepo, nil)
	_, err := uc.UpdateChat(ctx, userID, updateChat)

	assert.Error(t, err)
	assert.Equal(t, usecase.ErrOnlyOwnerCanModify, err)
}

func TestDeleteChat_NotOwner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	chatID := uuid.New()

	mockRepo.EXPECT().GetUserRoleInChat(ctx, userID, chatID).
		Return("member", nil) // Не owner

	uc := usecase.NewChatUsecase(mockRepo, nil)
	err := uc.DeleteChat(ctx, userID, chatID)

	assert.Error(t, err)
	assert.Equal(t, usecase.ErrOnlyOwnerCanModify, err)
}

func TestUpdateChat_DialogForbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	chatID := uuid.New()
	title := "New Title"
	updateChat := &model.UpdateChat{
		ID:    chatID,
		Title: &title,
	}

	mockRepo.EXPECT().GetChatByID(ctx, chatID).
		Return(&model.Chat{ID: chatID, Type: "dialog"}, nil)

	uc := usecase.NewChatUsecase(mockRepo, nil)
	_, err := uc.UpdateChat(ctx, userID, updateChat)

	assert.Error(t, err)
	assert.Equal(t, usecase.ErrDialogUpdateForbidden, err)
}

func TestCreateChat_WithAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	chatID := uuid.New()
	avatarURL := "path/to/avatar.png"
	createChat := &model.CreateChat{
		Title: "Group with avatar",
		Type:  "group",
	}

	// Мокаем создание чата с возвратом URL аватара
	mockRepo.EXPECT().
		CreateChat(ctx, createChat).
		Return(chatID, avatarURL, nil)

	// Мокаем добавление владельца
	mockRepo.EXPECT().
		AddUserToChatByID(ctx, userID, "owner", chatID).
		Return(nil)

	// Мокаем получение информации о чате
	mockRepo.EXPECT().
		GetUserRoleInChat(ctx, userID, chatID).
		Return("owner", nil)
	mockRepo.EXPECT().
		GetChatByID(ctx, chatID).
		Return(&model.Chat{
			ID:         chatID,
			Type:       "group",
			Title:      "Group with avatar",
			AvatarPath: &avatarURL,
		}, nil)
	mockRepo.EXPECT().
		GetUsersFromChat(ctx, chatID).
		Return([]model.UserInChat{}, nil)

	uc := usecase.NewChatUsecase(mockRepo, nil)
	info, err := uc.CreateChat(ctx, userID, createChat)

	assert.NoError(t, err)
	assert.Equal(t, avatarURL, *info.AvatarPath)
}

func TestUpdateChat_WithAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIChatRepo(ctrl)
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	chatID := uuid.New()
	newAvatarURL := "path/to/new_avatar.png"
	oldAvatarURL := "path/to/old_avatar.png"
	title := "Updated Title"
	updateChat := &model.UpdateChat{
		ID:    chatID,
		Title: &title,
	}

	mockRepo.EXPECT().GetChatByID(ctx, chatID).
		Return(&model.Chat{ID: chatID, Type: "group"}, nil)
	mockRepo.EXPECT().GetUserRoleInChat(ctx, userID, chatID).
		Return("owner", nil)
	mockRepo.EXPECT().UpdateChat(ctx, updateChat).
		Return(newAvatarURL, oldAvatarURL, nil)
	mockRepo.EXPECT().GetUserRoleInChat(ctx, userID, chatID).
		Return("owner", nil)
	mockRepo.EXPECT().GetChatByID(ctx, chatID).
		Return(&model.Chat{
			ID:         chatID,
			Type:       "group",
			Title:      title,
			AvatarPath: &newAvatarURL,
		}, nil)
	mockRepo.EXPECT().GetUsersFromChat(ctx, chatID).
		Return([]model.UserInChat{}, nil)

	uc := usecase.NewChatUsecase(mockRepo, nil)
	info, err := uc.UpdateChat(ctx, userID, updateChat)

	assert.NoError(t, err)
	assert.Equal(t, title, info.Title)
	assert.Equal(t, newAvatarURL, *info.AvatarPath)
}
