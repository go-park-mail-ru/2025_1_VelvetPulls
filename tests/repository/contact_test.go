package repository_test

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetContacts(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewContactRepo(db)
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "username", "avatar_path"}).
		AddRow(uuid.New(), "John", "Doe", "johndoe", "avatar.png").
		AddRow(uuid.New(), "Jane", "Smith", "janesmith", "avatar2.png")

	query := regexp.QuoteMeta(`
		SELECT u.id, u.first_name, u.last_name, u.username, u.avatar_path
		FROM public.contact c
		JOIN public.user u ON c.contact_id = u.id
		WHERE c.user_id = $1`)

	mock.ExpectQuery(query).WithArgs(userID).WillReturnRows(rows)

	contacts, err := repo.GetContacts(ctx, userID)
	require.NoError(t, err)
	require.Len(t, contacts, 2)
	assert.Equal(t, "johndoe", contacts[0].Username)
	assert.Equal(t, "janesmith", contacts[1].Username)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddContactByUsername_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewContactRepo(db)
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	userID := uuid.New()
	contactID := uuid.New()
	username := "friend"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM public.user WHERE username = $1")).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(contactID))

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO public.contact (user_id, contact_id) VALUES ($1, $2) ON CONFLICT DO NOTHING")).
		WithArgs(userID, contactID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AddContactByUsername(ctx, userID, username)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddContactByUsername_Self(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewContactRepo(db)
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	userID := uuid.New()
	username := "self"

	// Смокаем, что по username вернётся тот же ID
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM public.user WHERE username = $1")).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

	err = repo.AddContactByUsername(ctx, userID, username)
	assert.ErrorIs(t, err, repository.ErrSelfContact)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddContactByUsername_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewContactRepo(db)
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	userID := uuid.New()
	username := "notfound"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM public.user WHERE username = $1")).
		WithArgs(username).
		WillReturnError(sql.ErrNoRows)

	err = repo.AddContactByUsername(ctx, userID, username)
	assert.ErrorIs(t, err, repository.ErrUserNotFound)
}

func TestDeleteContactByUsername_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewContactRepo(db)
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	userID := uuid.New()
	contactID := uuid.New()
	username := "toDelete"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM public.user WHERE username = $1")).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(contactID))

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM public.contact WHERE user_id = $1 AND contact_id = $2")).
		WithArgs(userID, contactID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteContactByUsername(ctx, userID, username)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
