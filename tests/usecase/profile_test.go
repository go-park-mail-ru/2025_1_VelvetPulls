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

func TestGetUserProfileByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем мок-репозиторий пользователя
	mockUserRepo := mocks.NewMockIUserRepo(ctrl)
	userUC := usecase.NewUserUsecase(mockUserRepo)

	// Подготовка контекста с логгером
	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()

	// Подготовка тестовых данных (dummy user)
	firstName := "John"
	lastName := "Doe"
	username := "johndoe"
	phone := "1234567890"
	email := "john@example.com"
	avatarPath := "path/to/avatar.png"

	// Предполагаемая структура модели User (используем те же поля, что и в реализации usecase)
	mockUser := &model.User{
		FirstName:  &firstName,
		LastName:   &lastName,
		Username:   username,
		Phone:      phone,
		Email:      &email,
		AvatarPath: &avatarPath,
	}

	mockUserRepo.EXPECT().GetUserByID(ctx, userID).Return(mockUser, nil)

	profile, err := userUC.GetUserProfileByID(ctx, userID)
	require.NoError(t, err)
	require.NotNil(t, profile)

	assert.Equal(t, mockUser.FirstName, profile.FirstName)
	assert.Equal(t, mockUser.LastName, profile.LastName)
	assert.Equal(t, mockUser.Username, profile.Username)
	assert.Equal(t, mockUser.Phone, profile.Phone)
	assert.Equal(t, mockUser.Email, profile.Email)
	assert.Equal(t, mockUser.AvatarPath, profile.AvatarPath)
}

func TestGetUserProfileByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepo(ctrl)
	userUC := usecase.NewUserUsecase(mockUserRepo)

	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	username := "johndoe"

	// Тестовые данные
	firstName := "John"
	lastName := "Doe"
	phone := "1234567890"
	email := "john@example.com"
	avatarPath := "path/to/avatar.png"

	mockUser := &model.User{
		FirstName:  &firstName,
		LastName:   &lastName,
		Username:   username,
		Phone:      phone,
		Email:      &email,
		AvatarPath: &avatarPath,
	}

	mockUserRepo.EXPECT().GetUserByUsername(ctx, username).Return(mockUser, nil)

	profile, err := userUC.GetUserProfileByUsername(ctx, username)
	require.NoError(t, err)
	require.NotNil(t, profile)

	assert.Equal(t, mockUser.FirstName, profile.FirstName)
	assert.Equal(t, mockUser.LastName, profile.LastName)
	assert.Equal(t, mockUser.Username, profile.Username)
	assert.Equal(t, mockUser.Phone, profile.Phone)
	assert.Equal(t, mockUser.Email, profile.Email)
	assert.Equal(t, mockUser.AvatarPath, profile.AvatarPath)
}

func TestUpdateUserProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepo(ctrl)
	userUC := usecase.NewUserUsecase(mockUserRepo)

	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	newFirstName := "Jane"

	// Создаем валидный объект для обновления профиля: хотя бы одно поле должно быть задано
	updateProfile := &model.UpdateUserProfile{
		ID:        userID,
		FirstName: &newFirstName,
	}

	// Ожидается вызов UpdateUser, возвращающий пустые строки для нового и старого аватара и nil-ошибку
	mockUserRepo.EXPECT().UpdateUser(ctx, updateProfile).Return("", "", nil)

	// Вызов usecase
	err := userUC.UpdateUserProfile(ctx, updateProfile)
	require.NoError(t, err)
}

func TestUpdateUserProfile_InvalidUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepo(ctrl)
	userUC := usecase.NewUserUsecase(mockUserRepo)

	logger := zap.NewNop()
	ctx := context.WithValue(context.Background(), utils.LOGGER_ID_KEY, logger)

	userID := uuid.New()
	// Обновление профиля без указания ни одного поля должно вызвать ошибку валидации
	updateProfile := &model.UpdateUserProfile{
		ID: userID,
		// Все поля остаются nil
	}

	err := userUC.UpdateUserProfile(ctx, updateProfile)
	require.Error(t, err)
}
