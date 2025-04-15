package usecase_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	mocks "github.com/go-park-mail-ru/2025_1_VelvetPulls/tests/usecase/mock"
)

func TestGetUserContacts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContactRepo := mocks.NewMockIContactRepo(ctrl)
	contactUC := usecase.NewContactUsecase(mockContactRepo)

	// Формируем контекст с ненулевым логгером.
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.LOGGER_ID_KEY, zap.NewNop())

	userID := uuid.New()
	expectedContacts := []model.Contact{
		{ID: uuid.New(), FirstName: nil, LastName: nil, Username: "contact1", AvatarURL: nil},
		{ID: uuid.New(), FirstName: nil, LastName: nil, Username: "contact2", AvatarURL: nil},
	}

	// Ожидаем вызов метода GetContacts
	mockContactRepo.EXPECT().
		GetContacts(ctx, userID).
		Return(expectedContacts, nil)

	contacts, err := contactUC.GetUserContacts(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, expectedContacts, contacts)
}

func TestAddUserContact(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContactRepo := mocks.NewMockIContactRepo(ctrl)
	contactUC := usecase.NewContactUsecase(mockContactRepo)

	// Создаем контекст с логгером
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.LOGGER_ID_KEY, zap.NewNop())

	userID := uuid.New()
	contactUsername := "newcontact"

	// Ожидаем вызов метода AddContactByUsername
	mockContactRepo.EXPECT().
		AddContactByUsername(ctx, userID, contactUsername).
		Return(nil)

	err := contactUC.AddUserContact(ctx, userID, contactUsername)
	assert.NoError(t, err)
}

func TestRemoveUserContact(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContactRepo := mocks.NewMockIContactRepo(ctrl)
	contactUC := usecase.NewContactUsecase(mockContactRepo)

	// Контекст с логгером
	ctx := context.Background()
	ctx = context.WithValue(ctx, utils.LOGGER_ID_KEY, zap.NewNop())

	userID := uuid.New()
	contactUsername := "contactToRemove"

	// Ожидаем вызов метода DeleteContactByUsername
	mockContactRepo.EXPECT().
		DeleteContactByUsername(ctx, userID, contactUsername).
		Return(nil)

	err := contactUC.RemoveUserContact(ctx, userID, contactUsername)
	assert.NoError(t, err)
}
