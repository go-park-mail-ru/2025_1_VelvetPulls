package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	delivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/delivery/http"
	mocks "github.com/go-park-mail-ru/2025_1_VelvetPulls/tests/delivery/mock"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func addUserIDToContext(r *http.Request, userID uuid.UUID) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, utils.USER_ID_KEY, userID)
	ctx = context.WithValue(ctx, utils.LOGGER_ID_KEY, zap.NewNop())
	return r.WithContext(ctx)
}

func TestGetChats_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatUC := mocks.NewMockIChatUsecase(ctrl)
	mockSessUC := mocks.NewMockISessionUsecase(ctrl)

	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	expectedChats := []model.Chat{
		{
			ID:        uuid.New(),
			Title:     "Chat 1",
			Type:      "group",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
		{
			ID:        uuid.New(),
			Title:     "Chat 2",
			Type:      "dialog",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		},
	}

	// Ожидание вызова CheckLogin от мидлвари
	mockSessUC.EXPECT().
		CheckLogin(gomock.Any(), gomock.Any()).
		Return(userID.String(), nil)

	// Ожидание вызова GetChats
	mockChatUC.EXPECT().
		GetChats(gomock.Any(), userID).
		Return(expectedChats, nil)

	router := mux.NewRouter()
	delivery.NewChatController(router, mockChatUC, mockSessUC)

	req := httptest.NewRequest(http.MethodGet, "/chats", nil)

	// Добавляем заглушку куки (можно любую строку — мидлварь всё равно мокается)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "00000000-0000-0000-0000-000000000001",
	})
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp utils.JSONResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Status)
	assert.Empty(t, resp.Error)

	dataBytes, err := json.Marshal(resp.Data)
	assert.NoError(t, err)
	var chats []model.Chat
	err = json.Unmarshal(dataBytes, &chats)
	assert.NoError(t, err)
	assert.Len(t, chats, len(expectedChats))
}

func TestCreateChat_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatUC := mocks.NewMockIChatUsecase(ctrl)
	mockSessUC := mocks.NewMockISessionUsecase(ctrl)

	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	createReq := model.CreateChatRequest{
		Type:  "group",
		Title: "New Chat",
	}
	chatDataBytes, err := json.Marshal(createReq)
	assert.NoError(t, err)

	expectedChatInfo := &model.ChatInfo{
		ID:    uuid.New(),
		Title: createReq.Title,
	}
	mockSessUC.EXPECT().
		CheckLogin(gomock.Any(), gomock.Any()).
		Return(userID.String(), nil)

	mockChatUC.EXPECT().
		CreateChat(gomock.Any(), userID, gomock.AssignableToTypeOf(&model.CreateChat{})).
		DoAndReturn(func(ctx context.Context, uid uuid.UUID, chat *model.CreateChat) (*model.ChatInfo, error) {
			assert.Equal(t, createReq.Type, chat.Type)
			assert.Equal(t, createReq.Title, chat.Title)
			assert.Nil(t, chat.Avatar)
			return expectedChatInfo, nil
		})

	router := mux.NewRouter()
	delivery.NewChatController(router, mockChatUC, mockSessUC)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	partWriter, err := writer.CreateFormField("chat_data")
	assert.NoError(t, err)
	_, err = partWriter.Write(chatDataBytes)
	assert.NoError(t, err)
	err = writer.Close()
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/chat", body)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "00000000-0000-0000-0000-000000000001",
	})
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = addUserIDToContext(req, userID)
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp utils.JSONResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Status)
	assert.Empty(t, resp.Error)

	dataBytes, err := json.Marshal(resp.Data)
	assert.NoError(t, err)
	var chatInfo model.ChatInfo
	err = json.Unmarshal(dataBytes, &chatInfo)
	assert.NoError(t, err)
	assert.Equal(t, expectedChatInfo.Title, chatInfo.Title)
}

func TestGetChat_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatUC := mocks.NewMockIChatUsecase(ctrl)
	mockSessUC := mocks.NewMockISessionUsecase(ctrl)

	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	chatID := uuid.New()

	expectedChatInfo := &model.ChatInfo{
		ID:    chatID,
		Title: "Chat Info",
	}
	mockSessUC.EXPECT().
		CheckLogin(gomock.Any(), gomock.Any()).
		Return(userID.String(), nil)

	mockChatUC.EXPECT().
		GetChatInfo(gomock.Any(), userID, chatID).
		Return(expectedChatInfo, nil)

	router := mux.NewRouter()
	delivery.NewChatController(router, mockChatUC, mockSessUC)

	url := "/chat/" + chatID.String()
	req := httptest.NewRequest(http.MethodGet, url, nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "00000000-0000-0000-0000-000000000001",
	})
	req = addUserIDToContext(req, userID)
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp utils.JSONResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Status)
	assert.Empty(t, resp.Error)

	dataBytes, err := json.Marshal(resp.Data)
	assert.NoError(t, err)
	var chatInfo model.ChatInfo
	err = json.Unmarshal(dataBytes, &chatInfo)
	assert.NoError(t, err)
	assert.Equal(t, expectedChatInfo.Title, chatInfo.Title)
}

func TestUpdateChat_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChatUC := mocks.NewMockIChatUsecase(ctrl)
	mockSessUC := mocks.NewMockISessionUsecase(ctrl)

	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	chatID := uuid.New()
	updateTitle := "Updated Chat Title"

	updateReq := model.UpdateChat{
		ID:    chatID,
		Title: &updateTitle,
	}
	updateDataBytes, err := json.Marshal(updateReq)
	assert.NoError(t, err)

	expectedChatInfo := &model.ChatInfo{
		ID:    chatID,
		Title: updateTitle,
	}

	mockSessUC.EXPECT().
		CheckLogin(gomock.Any(), gomock.Any()).
		Return(userID.String(), nil)

	mockChatUC.EXPECT().
		UpdateChat(gomock.Any(), userID, gomock.AssignableToTypeOf(&model.UpdateChat{})).
		DoAndReturn(func(ctx context.Context, uid uuid.UUID, upd *model.UpdateChat) (*model.ChatInfo, error) {
			assert.Equal(t, chatID, upd.ID)
			assert.NotNil(t, upd.Title)
			return expectedChatInfo, nil
		})

	router := mux.NewRouter()
	delivery.NewChatController(router, mockChatUC, mockSessUC)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	partWriter, err := writer.CreateFormField("chat_data")
	assert.NoError(t, err)
	_, err = partWriter.Write(updateDataBytes)
	assert.NoError(t, err)
	err = writer.Close()
	assert.NoError(t, err)

	url := "/chat/" + chatID.String()
	req := httptest.NewRequest(http.MethodPut, url, body)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "00000000-0000-0000-0000-000000000001",
	})
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req = addUserIDToContext(req, userID)
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp utils.JSONResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Status)
	assert.Empty(t, resp.Error)

	dataBytes, err := json.Marshal(resp.Data)
	assert.NoError(t, err)
	var chatInfo model.ChatInfo
	err = json.Unmarshal(dataBytes, &chatInfo)
	assert.NoError(t, err)
	assert.Equal(t, expectedChatInfo.Title, chatInfo.Title)
}
