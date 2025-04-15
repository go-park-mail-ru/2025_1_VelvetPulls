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

	usecase := usecase.NewChatUsecase(mockRepo)
	chat, err := usecase.GetChats(ctx, chatID)

	assert.NoError(t, err)
	assert.Equal(t, []model.Chat{*expectedChat}, chat)
}

func TestGetChatInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatRepo := mocks.NewMockIChatRepo(ctrl)
	chatUC := usecase.NewChatUsecase(mockChatRepo)

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

	uc := usecase.NewChatUsecase(mockRepo)
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

	uc := usecase.NewChatUsecase(mockRepo)
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

	usecase := usecase.NewChatUsecase(mockRepo)
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

	usecase := usecase.NewChatUsecase(mockRepo)
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

	usecase := usecase.NewChatUsecase(mockRepo)
	res, err := usecase.DeleteUserFromChat(ctx, userID, toDelete, chatID)

	assert.NoError(t, err)
	assert.Equal(t, toDelete, res.DeletedUsers)
}
