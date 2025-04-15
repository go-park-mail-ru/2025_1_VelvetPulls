package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	delivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/delivery/http"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	mocks "github.com/go-park-mail-ru/2025_1_VelvetPulls/tests/delivery/mock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestGetMessageHistory_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageUC := mocks.NewMockIMessageUsecase(ctrl)
	mockSessionUC := mocks.NewMockISessionUsecase(ctrl)

	// Use fixed UUIDs for predictable comparisons
	userID := uuid.MustParse("5caed5b9-1315-4113-b1d1-ea34b3d771e5")
	chatID := uuid.MustParse("183bb09c-8e4a-495d-92c0-77667fd2f585")
	msg1ID := uuid.MustParse("aae03890-a767-4680-b293-44f9f825ad7f")
	msg2ID := uuid.MustParse("3193fe6e-2db0-46dd-bd92-4da04c242125")

	now := time.Now()
	expectedMessages := []model.Message{
		{
			ID:         msg1ID,
			ChatID:     chatID,
			UserID:     userID,
			Body:       "Test message 1",
			SentAt:     now,
			IsRedacted: false,
		},
		{
			ID:         msg2ID,
			ChatID:     chatID,
			UserID:     userID,
			Body:       "Test message 2",
			SentAt:     now.Add(time.Second),
			IsRedacted: false,
		},
	}

	mockSessionUC.EXPECT().
		CheckLogin(gomock.Any(), "valid-token").
		Return(userID.String(), nil)

	mockMessageUC.EXPECT().
		GetChatMessages(gomock.Any(), userID, chatID).
		Return(expectedMessages, nil)

	router := mux.NewRouter()
	delivery.NewMessageController(router, mockMessageUC, mockSessionUC)

	req := httptest.NewRequest(http.MethodGet, "/chat/"+chatID.String()+"/messages", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "valid-token",
	})
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp utils.JSONResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Status)

	var actualMessages []model.Message
	jsonData, err := json.Marshal(resp.Data)
	assert.NoError(t, err)
	err = json.Unmarshal(jsonData, &actualMessages)
	assert.NoError(t, err)

	// Compare messages ignoring exact time values
	assert.Equal(t, len(expectedMessages), len(actualMessages))
	for i := range expectedMessages {
		assert.Equal(t, expectedMessages[i].ID, actualMessages[i].ID)
		assert.Equal(t, expectedMessages[i].ChatID, actualMessages[i].ChatID)
		assert.Equal(t, expectedMessages[i].UserID, actualMessages[i].UserID)
		assert.Equal(t, expectedMessages[i].Body, actualMessages[i].Body)
		assert.Equal(t, expectedMessages[i].IsRedacted, actualMessages[i].IsRedacted)
		// Don't compare SentAt exactly, just check it's not zero
		assert.False(t, actualMessages[i].SentAt.IsZero())
	}
}

func TestGetMessageHistory_InvalidChatID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageUC := mocks.NewMockIMessageUsecase(ctrl)
	mockSessionUC := mocks.NewMockISessionUsecase(ctrl)

	userID := uuid.New()

	mockSessionUC.EXPECT().
		CheckLogin(gomock.Any(), "valid-token").
		Return(userID.String(), nil)

	router := mux.NewRouter()
	delivery.NewMessageController(router, mockMessageUC, mockSessionUC)

	req := httptest.NewRequest(http.MethodGet, "/chat/invalid/messages", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "valid-token",
	})
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSendMessage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageUC := mocks.NewMockIMessageUsecase(ctrl)
	mockSessionUC := mocks.NewMockISessionUsecase(ctrl)

	userID := uuid.New()
	chatID := uuid.New()
	messageInput := model.MessageInput{
		Message: "Test message",
	}

	mockSessionUC.EXPECT().
		CheckLogin(gomock.Any(), "valid-token").
		Return(userID.String(), nil)

	mockMessageUC.EXPECT().
		SendMessage(gomock.Any(), &messageInput, userID, chatID).
		Return(nil)

	router := mux.NewRouter()
	delivery.NewMessageController(router, mockMessageUC, mockSessionUC)

	body, err := json.Marshal(messageInput)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/chat/"+chatID.String()+"/messages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "valid-token",
	})
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp utils.JSONResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Status)
	assert.Equal(t, "message send successful", resp.Data)
}

func TestSendMessage_InvalidBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageUC := mocks.NewMockIMessageUsecase(ctrl)
	mockSessionUC := mocks.NewMockISessionUsecase(ctrl)

	userID := uuid.New()
	chatID := uuid.New()

	mockSessionUC.EXPECT().
		CheckLogin(gomock.Any(), "valid-token").
		Return(userID.String(), nil)

	router := mux.NewRouter()
	delivery.NewMessageController(router, mockMessageUC, mockSessionUC)

	req := httptest.NewRequest(http.MethodPost, "/chat/"+chatID.String()+"/messages", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "valid-token",
	})
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestSendMessage_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessageUC := mocks.NewMockIMessageUsecase(ctrl)
	mockSessionUC := mocks.NewMockISessionUsecase(ctrl)

	chatID := uuid.New()
	messageInput := model.MessageInput{
		Message: "Test message",
	}

	router := mux.NewRouter()
	delivery.NewMessageController(router, mockMessageUC, mockSessionUC)

	body, err := json.Marshal(messageInput)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/chat/"+chatID.String()+"/messages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
