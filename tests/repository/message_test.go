package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateMessage(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewMessageRepo(db)

	messageID := uuid.New()
	userID := uuid.New()
	chatID := uuid.New()
	sentAt := time.Now()
	username := "test_user"
	avatarPath := "/avatar.jpg"
	body := "hello"

	msg := &model.Message{
		UserID: userID,
		ChatID: chatID,
		Body:   body,
	}

	// Expect INSERT
	mock.ExpectQuery(`INSERT INTO message`).
		WithArgs(userID, chatID, body).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(messageID))

	// Expect SELECT after insert (getMessage)
	mock.ExpectQuery(`SELECT(.|\n)*FROM message m(.|\n)*WHERE m.id = \$1`).
		WithArgs(messageID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "parent_message_id", "chat_id", "user_id", "body", "sent_at", "is_redacted", "username", "avatar_path",
		}).AddRow(messageID, nil, chatID, userID, body, sentAt, false, username, avatarPath))

	result, err := repo.CreateMessage(context.Background(), msg)
	require.NoError(t, err)
	require.Equal(t, messageID, result.ID)
	require.Equal(t, username, result.Username)
	if assert.NotNil(t, result.AvatarPath) {
		require.Equal(t, "/avatar.jpg", *result.AvatarPath)
	}
	require.Equal(t, body, result.Body)
}

func TestGetMessages(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewMessageRepo(db)

	chatID := uuid.New()
	messageID := uuid.New()
	userID := uuid.New()
	sentAt := time.Now()
	body := "hello"
	username := "test_user"
	avatarPath := "/avatar.jpg"

	rows := sqlmock.NewRows([]string{
		"id", "parent_message_id", "chat_id", "user_id", "body", "sent_at", "is_redacted", "username", "avatar_path",
	}).AddRow(messageID, nil, chatID, userID, body, sentAt, false, username, avatarPath)

	mock.ExpectQuery(`SELECT(.|\n)*FROM message m(.|\n)*WHERE m.chat_id = \$1`).
		WithArgs(chatID).
		WillReturnRows(rows)

	result, err := repo.GetMessages(context.Background(), chatID)
	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, messageID, result[0].ID)
	require.Equal(t, body, result[0].Body)
	require.Equal(t, username, result[0].Username)
}
