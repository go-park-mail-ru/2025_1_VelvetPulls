package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	mocks "github.com/go-park-mail-ru/2025_1_VelvetPulls/tests/usecase/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepo(ctrl)
	mockSessionRepo := mocks.NewMockISessionRepo(ctrl)
	authUC := usecase.NewAuthUsecase(mockUserRepo, mockSessionRepo)

	// Создаем контекст и сразу устанавливаем в него ненулевой логгер
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.LOGGER_ID_KEY, zap.NewNop())

	creds := model.RegisterCredentials{
		Username:        "testuser",
		Password:        "Password123!",
		ConfirmPassword: "Password123!",
		Phone:           "1234567890",
	}

	// Ожидаем, что пользователь не найден по username и по телефону
	mockUserRepo.EXPECT().GetUserByUsername(ctx, creds.Username).
		Return(nil, errors.New("not found"))
	mockUserRepo.EXPECT().GetUserByPhone(ctx, creds.Phone).
		Return(nil, errors.New("not found"))

	// Затем ожидаем создание пользователя и последующее создание сессии
	mockUserRepo.EXPECT().CreateUser(ctx, gomock.Any()).
		Return("user-id-123", nil)
	mockSessionRepo.EXPECT().CreateSession(ctx, "user-id-123").
		Return("session-id-abc", nil)

	sessionID, err := authUC.RegisterUser(ctx, creds)
	assert.NoError(t, err)
	assert.Equal(t, "session-id-abc", sessionID)
}

func TestRegisterUser_UsernameTaken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepo(ctrl)
	mockSessionRepo := mocks.NewMockISessionRepo(ctrl)
	authUC := usecase.NewAuthUsecase(mockUserRepo, mockSessionRepo)

	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.LOGGER_ID_KEY, zap.NewNop())

	creds := model.RegisterCredentials{
		Username:        "existinguser",
		Password:        "Password123!",
		ConfirmPassword: "Password123!",
		Phone:           "1234567890",
	}

	mockUserRepo.EXPECT().GetUserByUsername(ctx, creds.Username).Return(&model.User{}, nil)

	_, err := authUC.RegisterUser(ctx, creds)
	assert.ErrorIs(t, err, usecase.ErrUsernameIsTaken)
}

func TestLoginUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepo(ctrl)
	mockSessionRepo := mocks.NewMockISessionRepo(ctrl)
	authUC := usecase.NewAuthUsecase(mockUserRepo, mockSessionRepo)

	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.LOGGER_ID_KEY, zap.NewNop())
	hashedPass, _ := utils.HashAndSalt("Password123!")

	creds := model.LoginCredentials{
		Username: "testuser",
		Password: "Password123!",
	}

	mockUserRepo.EXPECT().GetUserByUsername(ctx, creds.Username).Return(&model.User{
		ID:       uuid.New(),
		Username: "testuser",
		Password: hashedPass,
	}, nil)

	mockSessionRepo.EXPECT().CreateSession(ctx, gomock.Any()).Return("session-id-xyz", nil)

	sessionID, err := authUC.LoginUser(ctx, creds)
	assert.NoError(t, err)
	assert.Equal(t, "session-id-xyz", sessionID)
}

func TestLogoutUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockIUserRepo(ctrl)
	mockSessionRepo := mocks.NewMockISessionRepo(ctrl)
	authUC := usecase.NewAuthUsecase(mockUserRepo, mockSessionRepo)

	// Создаем контекст с ненулевым логгером, чтобы utils.GetLoggerFromCtx(ctx) не паниковал.
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.LOGGER_ID_KEY, zap.NewNop())

	sessionID := "session-123"

	// Ожидаем вызов метода DeleteSession из session репозитория
	mockSessionRepo.EXPECT().
		DeleteSession(ctx, sessionID).
		Return(nil)

	err := authUC.LogoutUser(ctx, sessionID)
	assert.NoError(t, err)
}
