package repository_test

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetUserIDByToken_Success(t *testing.T) {
	db, mock := redismock.NewClientMock()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())
	sessionID := "test-session-id"
	expectedUserID := "123e4567-e89b-12d3-a456-426614174000"

	mock.ExpectGet(sessionID).SetVal(expectedUserID)

	repo := repository.NewSessionRepo(db)
	userID, err := repo.GetUserIDByToken(ctx, sessionID)

	require.NoError(t, err)
	require.Equal(t, expectedUserID, userID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateSession_Success(t *testing.T) {
	db, mock := redismock.NewClientMock()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())
	userID := "123e4567-e89b-12d3-a456-426614174000"

	mock.Regexp().ExpectSet(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, userID, config.CookieDuration).SetVal("OK")

	repo := repository.NewSessionRepo(db)
	sessionID, err := repo.CreateSession(ctx, userID)

	require.NoError(t, err)
	require.NotEmpty(t, sessionID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteSession_Success(t *testing.T) {
	db, mock := redismock.NewClientMock()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())
	sessionID := "test-session-id"

	mock.ExpectExists(sessionID).SetVal(1)
	mock.ExpectDel(sessionID).SetVal(1)

	repo := repository.NewSessionRepo(db)
	err := repo.DeleteSession(ctx, sessionID)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
