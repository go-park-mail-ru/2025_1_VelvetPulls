package repository_test

import (
	"context"
	"database/sql"
	"errors"
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

func TestCreateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	newUser := &model.User{
		Username: "newuser",
		Phone:    "9876543210",
		Password: "hashedpass",
	}
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	// Ожидаем успешное выполнение запроса INSERT и возврат нового id.
	expectedID := uuid.New()
	mock.ExpectQuery(`INSERT INTO public.user \(username, phone, password\) VALUES \(\$1, \$2, \$3\) RETURNING id`).
		WithArgs(newUser.Username, newUser.Phone, newUser.Password).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	userID, err := repo.CreateUser(ctx, newUser)
	require.NoError(t, err)
	require.Equal(t, expectedID.String(), userID)
}

// -----------------------------
// Остальные тесты (GetUserByEmail, GetUserByPhone, GetUserByID, CreateUser (ошибки), UpdateUser и т.д.)
// уже присутствуют в вашем файле и ниже приведены для справки.
// -----------------------------

func TestGetUserByEmail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	email := "test@mail.com"
	expectedID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"id", "avatar_path", "first_name", "last_name", "username",
		"phone", "email", "password", "created_at", "updated_at",
	}).AddRow(
		expectedID, "/avatar.jpg", "John", "Doe", "testuser",
		"1234567890", email, "hashedpass", time.Now(), time.Now(),
	)

	mock.ExpectQuery(`SELECT .* FROM public.user WHERE email = \$1`).
		WithArgs(email).
		WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	user, err := repo.GetUserByEmail(ctx, email)
	require.NoError(t, err)
	require.Equal(t, email, *user.Email)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	email := "missing@mail.com"
	mock.ExpectQuery(`SELECT .* FROM public.user WHERE email = \$1`).
		WithArgs(email).
		WillReturnError(sql.ErrNoRows)

	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	user, err := repo.GetUserByEmail(ctx, email)
	require.Nil(t, user)
	require.ErrorIs(t, err, repository.ErrUserNotFound)
}

func TestGetUserByPhone_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	phone := "1234567890"
	expectedID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"id", "avatar_path", "first_name", "last_name", "username",
		"phone", "email", "password", "created_at", "updated_at",
	}).AddRow(
		expectedID, "/avatar.jpg", "John", "Doe", "testuser",
		phone, "test@mail.com", "hashedpass", time.Now(), time.Now(),
	)

	mock.ExpectQuery(`SELECT .* FROM public.user WHERE phone = \$1`).
		WithArgs(phone).
		WillReturnRows(rows)

	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	user, err := repo.GetUserByPhone(ctx, phone)
	require.NoError(t, err)
	require.Equal(t, phone, user.Phone)
}

func TestGetUserByPhone_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	phone := "0000000000"
	mock.ExpectQuery(`SELECT .* FROM public.user WHERE phone = \$1`).
		WithArgs(phone).
		WillReturnError(sql.ErrNoRows)

	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	user, err := repo.GetUserByPhone(ctx, phone)
	require.Nil(t, user)
	require.ErrorIs(t, err, repository.ErrUserNotFound)
}

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

func TestCreateUser_DuplicateUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	user := &model.User{
		Username: "existinguser",
		Phone:    "1234567890",
		Password: "hashedpass",
	}
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	mock.ExpectQuery(`INSERT INTO public.user .* RETURNING id`).
		WithArgs(user.Username, user.Phone, user.Password).
		WillReturnError(errors.New("duplicate key value violates unique constraint \"user_username_key\""))

	_, err = repo.CreateUser(ctx, user)
	require.ErrorIs(t, err, repository.ErrRecordAlreadyExists)
}

func TestCreateUser_DuplicatePhone(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	user := &model.User{
		Username: "newuser",
		Phone:    "existingphone",
		Password: "hashedpass",
	}
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	mock.ExpectQuery(`INSERT INTO public.user .* RETURNING id`).
		WithArgs(user.Username, user.Phone, user.Password).
		WillReturnError(errors.New("duplicate key value violates unique constraint \"user_phone_key\""))

	_, err = repo.CreateUser(ctx, user)
	require.ErrorIs(t, err, repository.ErrRecordAlreadyExists)
}

func TestCreateUser_EmptyFields(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, zap.NewNop())

	testCases := []struct {
		name     string
		user     *model.User
		expected error
	}{
		{
			name:     "Nil user",
			user:     nil,
			expected: repository.ErrEmptyField,
		},
		{
			name:     "Empty username",
			user:     &model.User{Phone: "123", Password: "pass"},
			expected: repository.ErrEmptyField,
		},
		{
			name:     "Empty phone",
			user:     &model.User{Username: "user", Password: "pass"},
			expected: repository.ErrEmptyField,
		},
		{
			name:     "Empty password",
			user:     &model.User{Username: "user", Phone: "123"},
			expected: repository.ErrEmptyField,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := repo.CreateUser(ctx, tc.user)
			require.ErrorIs(t, err, tc.expected)
		})
	}
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
