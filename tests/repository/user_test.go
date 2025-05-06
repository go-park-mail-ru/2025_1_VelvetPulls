package repository_test

import (
	"context"
	"database/sql"
	"mime/multipart"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Test GetUserByEmail, GetUserByPhone, GetUserByID, CreateUser, UpdateUser уже имеются – ниже добавлены недостающие тесты.

// -----------------------------
// Новые тесты для GetUserByUsername
// -----------------------------

func TestGetUserByUsername_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	username := "testuser"
	expectedID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"id", "avatar_path", "first_name", "last_name", "username",
		"phone", "email", "password", "created_at", "updated_at",
	}).AddRow(
		expectedID, "/avatar.jpg", "John", "Doe", username,
		"1234567890", "test@mail.com", "hashedpass", time.Now(), time.Now(),
	)

	mock.ExpectQuery(`SELECT .* FROM public.user WHERE username = \$1`).
		WithArgs(username).
		WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	user, err := repo.GetUserByUsername(ctx, username)
	require.NoError(t, err)
	require.Equal(t, username, user.Username)
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	username := "notfound"
	mock.ExpectQuery(`SELECT .* FROM public.user WHERE username = \$1`).
		WithArgs(username).
		WillReturnError(sql.ErrNoRows)

	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	user, err := repo.GetUserByUsername(ctx, username)
	require.Nil(t, user)
	require.ErrorIs(t, err, repository.ErrUserNotFound)
}

// -----------------------------
// Тест для успешного создания пользователя (CreateUser)
// -----------------------------

func TestGetUserByID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	userID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"id", "avatar_path", "first_name", "last_name", "username",
		"phone", "email", "password", "created_at", "updated_at",
	}).AddRow(
		userID, "/avatar.jpg", "John", "Doe", "testuser",
		"1234567890", "test@mail.com", "hashedpass", time.Now(), time.Now(),
	)

	mock.ExpectQuery(`SELECT .* FROM public.user WHERE id = \$1`).
		WithArgs(userID).
		WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	user, err := repo.GetUserByID(ctx, userID)
	require.NoError(t, err)
	require.Equal(t, userID, user.ID)
}

func TestGetUserByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	userID := uuid.New()
	mock.ExpectQuery(`SELECT .* FROM public.user WHERE id = \$1`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	user, err := repo.GetUserByID(ctx, userID)
	require.Nil(t, user)
	require.ErrorIs(t, err, repository.ErrUserNotFound)
}

func TestGetUserByID_InvalidUUID(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	user, err := repo.GetUserByID(ctx, uuid.Nil)
	require.Nil(t, user)
	require.ErrorIs(t, err, repository.ErrInvalidUUID)
}

func TestUpdateUser_WithAvatar(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	userID := uuid.New()
	oldAvatar := "/old/avatar.jpg"
	newUsername := "newname"

	profile := &model.UpdateUserProfile{
		ID:       userID,
		Username: &newUsername,
	}

	// 1. Ожидаем начало транзакции
	mock.ExpectBegin()

	// 2. Ожидаем запрос на блокировку строки
	mock.ExpectQuery(`SELECT avatar_path FROM public.user WHERE id = \$1 FOR UPDATE`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"avatar_path"}).AddRow(oldAvatar))

	// 3. Ожидаем UPDATE запрос
	// Поскольку обновляются не все поля, порядковый номер аргументов может меняться.
	mock.ExpectExec(`UPDATE public.user SET (.+) WHERE id = \$4`).
		WithArgs(
			sqlmock.AnyArg(), // новый avatar_path
			newUsername,      // username
			sqlmock.AnyArg(), // updated_at
			userID,           // id
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// 4. Ожидаем коммит транзакции
	mock.ExpectCommit()

	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	// Для теста достаточно задать ненулевое значение для profile.Avatar.
	// Заметим, что тип profile.Avatar — это указатель на multipart.File.
	profile.Avatar = new(multipart.File)

	newAvatar, oldAvatarPath, err := repo.UpdateUser(ctx, profile)
	require.NoError(t, err)
	require.NotEmpty(t, newAvatar)
	require.Equal(t, oldAvatar, oldAvatarPath)
}

func TestUpdateUser_NoChanges(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	userID := uuid.New()
	profile := &model.UpdateUserProfile{
		ID: userID,
	}

	// Если нет изменений (нет ни одного ненулевого поля) транзакция откатывается.
	mock.ExpectBegin()
	mock.ExpectRollback()

	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	_, _, err = repo.UpdateUser(ctx, profile)
	require.NoError(t, err)
}

func TestUpdateUser_EmptyField(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	userID := uuid.New()
	empty := ""
	profile := &model.UpdateUserProfile{
		ID:        userID,
		FirstName: &empty,
	}

	// Ожидаем начало транзакции и затем откат.
	mock.ExpectBegin()
	mock.ExpectRollback()

	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	_, _, err = repo.UpdateUser(ctx, profile)
	require.ErrorIs(t, err, repository.ErrEmptyField)
}

func TestUpdateUser_InvalidUUID(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	_, _, err = repo.UpdateUser(ctx, &model.UpdateUserProfile{ID: uuid.Nil})
	require.ErrorIs(t, err, repository.ErrInvalidUUID)
}
