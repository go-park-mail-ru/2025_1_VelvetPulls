package repository_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetChatByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewChatRepo(db)

	chatID := uuid.New()
	avatarPath := "path/to/avatar.png"
	expectedChat := model.Chat{
		ID:         chatID,
		AvatarPath: &avatarPath,
		Type:       "group",
		Title:      "Test Chat",
	}

	rows := sqlmock.NewRows([]string{"id", "avatar_path", "type", "title"}).
		AddRow(expectedChat.ID, expectedChat.AvatarPath, expectedChat.Type, expectedChat.Title)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, avatar_path, type, title FROM chat WHERE id = $1")).
		WithArgs(chatID).
		WillReturnRows(rows)

	ctx := context.Background()
	chat, err := repo.GetChatByID(ctx, chatID)
	require.NoError(t, err)
	require.NotNil(t, chat)
	assert.Equal(t, expectedChat.ID, chat.ID)
	assert.Equal(t, expectedChat.Title, chat.Title)
	assert.Equal(t, expectedChat.Type, chat.Type)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetChats(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewChatRepo(db)
	userID := uuid.New()

	chat1 := model.Chat{
		ID:         uuid.New(),
		AvatarPath: nil,
		Type:       "group",
		Title:      "Chat 1",
	}
	chat2 := model.Chat{
		ID:         uuid.New(),
		AvatarPath: nil,
		Type:       "dialog",
		Title:      "Chat 2",
	}

	rows := sqlmock.NewRows([]string{
		"c.id", "c.avatar_path", "c.type", "c.title", "uc.send_notifications",
		"m.id", "m.user_id", "m.body", "m.sent_at", "count_users",
	}).
		AddRow(chat1.ID, chat1.AvatarPath, chat1.Type, chat1.Title, true,
			nil, nil, nil, nil, 2).
		AddRow(chat2.ID, chat2.AvatarPath, chat2.Type, chat2.Title, false,
			nil, nil, nil, nil, 1)

	mock.ExpectQuery("(?s)SELECT c\\.id.*FROM chat c.*WHERE uc\\.user_id = \\$1.*ORDER BY m\\.sent_at DESC NULLS LAST").
		WithArgs(userID).
		WillReturnRows(rows)

	ctx := context.Background()
	chats, lastChatID, err := repo.GetChats(ctx, userID)
	require.NoError(t, err)
	require.Len(t, chats, 2)
	assert.Equal(t, chat2.ID, lastChatID) // предполагается, что он будет последним, если сортировка работает

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCreateChat_NoAvatar(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewChatRepo(db)
	ctx := context.Background()

	createChat := &model.CreateChatRequest{
		Type:  "group",
		Title: "New Group",
	}

	query := regexp.QuoteMeta(`
        INSERT INTO chat (type, title, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())
        RETURNING id
    `)

	chatID := uuid.New()
	mock.ExpectQuery(query).
		WithArgs(createChat.Type, createChat.Title).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(chatID))

	logger := zap.NewNop()
	ctx = context.WithValue(ctx, utils.LOGGER_ID_KEY, logger)

	retID, err := repo.CreateChat(ctx, createChat)
	require.NoError(t, err)
	assert.Equal(t, chatID, retID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateChat(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewChatRepo(db)
	ctx := context.Background()
	logger := zap.NewNop()
	ctx = context.WithValue(ctx, utils.LOGGER_ID_KEY, logger)

	chatID := uuid.New()
	title := "Updated Title"
	updateChat := &model.UpdateChat{
		ID:    chatID,
		Title: &title,
	}

	mock.ExpectBegin()

	query := regexp.QuoteMeta(`UPDATE chat SET title = $1, updated_at = $2 WHERE id = $3`)
	mock.ExpectExec(query).
		WithArgs(title, sqlmock.AnyArg(), chatID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	avatarNewURL, avatarOldURL, err := repo.UpdateChat(ctx, updateChat)
	require.NoError(t, err)
	assert.Equal(t, "", avatarNewURL)
	assert.Equal(t, "", avatarOldURL)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteChat(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewChatRepo(db)
	ctx := context.Background()

	chatID := uuid.New()
	query := regexp.QuoteMeta("DELETE FROM chat WHERE id = $1")
	mock.ExpectExec(query).
		WithArgs(chatID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteChat(ctx, chatID)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestAddUserToChatByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewChatRepo(db)
	ctx := context.Background()
	userID := uuid.New()
	chatID := uuid.New()
	userRole := "owner"

	query := regexp.QuoteMeta("INSERT INTO user_chat (user_id, chat_id, user_role, joined_at) VALUES ($1, $2, $3, NOW())")
	mock.ExpectExec(query).
		WithArgs(userID, chatID, userRole).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AddUserToChatByID(ctx, userID, userRole, chatID)
	assert.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserRoleInChat_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewChatRepo(db)
	ctx := context.Background()
	userID := uuid.New()
	chatID := uuid.New()

	query := regexp.QuoteMeta("SELECT user_role FROM user_chat WHERE user_id = $1 AND chat_id = $2")
	mock.ExpectQuery(query).
		WithArgs(userID, chatID).
		WillReturnError(sql.ErrNoRows)

	role, err := repo.GetUserRoleInChat(ctx, userID, chatID)
	require.NoError(t, err)
	// Если строк нет, согласно реализации, возвращается пустая строка.
	assert.Equal(t, "", role)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUsersFromChat(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewChatRepo(db)
	ctx := context.Background()
	chatID := uuid.New()

	// Подготавливаем один ряд для одного пользователя.
	// Столбцы: u.id, u.username, u.name, u.avatar_path, uc.user_role.
	username := "testuser"
	name := "John"
	avatarPath := "avatar.png"
	role := "member"

	rows := sqlmock.NewRows([]string{"id", "username", "name", "avatar_path", "user_role"}).
		AddRow(uuid.New(), username, name, avatarPath, role)

	query := regexp.QuoteMeta(`SELECT u.id, u.username, u.name, u.avatar_path, uc.user_role 
	FROM public.user u JOIN user_chat uc ON u.id = uc.user_id WHERE uc.chat_id = $1`)

	mock.ExpectQuery(query).
		WithArgs(chatID).
		WillReturnRows(rows)

	users, err := repo.GetUsersFromChat(ctx, chatID)
	require.NoError(t, err)
	require.Len(t, users, 1)

	assert.Equal(t, username, users[0].Username)
	assert.Equal(t, name, *users[0].Name)
	assert.Equal(t, avatarPath, *users[0].AvatarPath)
	assert.Equal(t, role, *users[0].Role)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetChatByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewChatRepo(db)
	chatID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, avatar_path, type, title FROM chat WHERE id = $1")).
		WithArgs(chatID).
		WillReturnError(sql.ErrNoRows)

	ctx := context.Background()
	_, err = repo.GetChatByID(ctx, chatID)
	assert.Error(t, err)
	assert.Equal(t, repository.ErrChatNotFound, err)
}

func TestAddUserToChatByUsername_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewChatRepo(db)
	ctx := context.Background()
	username := "nonexistent"
	chatID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM public.user WHERE username = $1")).
		WithArgs(username).
		WillReturnError(sql.ErrNoRows)

	err = repo.AddUserToChatByUsername(ctx, username, "member", chatID)
	assert.Error(t, err)
}

func TestRemoveUserFromChatByUsername_UserNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewChatRepo(db)
	ctx := context.Background()
	username := "nonexistent"
	chatID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM public.user WHERE username = $1")).
		WithArgs(username).
		WillReturnError(sql.ErrNoRows)

	err = repo.RemoveUserFromChatByUsername(ctx, username, chatID)
	assert.Error(t, err)
}
