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

	delivery "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/delivery/http"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/proto"
	mocks "github.com/go-park-mail-ru/2025_1_VelvetPulls/tests/delivery/mock"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
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
	mockSessionClient := mocks.NewMockSessionServiceClient(ctrl)

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

	// Настраиваем ожидания для gRPC клиента
	mockSessionClient.EXPECT().
		CheckLogin(gomock.Any(), &authpb.CheckLoginRequest{SessionId: "00000000-0000-0000-0000-000000000001"}, gomock.Any()).
		Return(&authpb.CheckLoginResponse{UserId: userID.String()}, nil)

	// Настраиваем ожидания для usecase
	mockChatUC.EXPECT().
		GetChats(gomock.Any(), userID).
		Return(expectedChats, nil)

	router := mux.NewRouter()
	delivery.NewChatController(router, mockChatUC, mockSessionClient)

	req := httptest.NewRequest(http.MethodGet, "/chats", nil)
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
	mockSessionClient := mocks.NewMockSessionServiceClient(ctrl)

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

	// Настраиваем ожидания для gRPC клиента
	mockSessionClient.EXPECT().
		CheckLogin(gomock.Any(), &authpb.CheckLoginRequest{SessionId: "00000000-0000-0000-0000-000000000001"}, gomock.Any()).
		Return(&authpb.CheckLoginResponse{UserId: userID.String()}, nil)

	mockChatUC.EXPECT().
		CreateChat(gomock.Any(), userID, gomock.AssignableToTypeOf(&model.CreateChat{})).
		DoAndReturn(func(ctx context.Context, uid uuid.UUID, chat *model.CreateChat) (*model.ChatInfo, error) {
			assert.Equal(t, createReq.Type, chat.Type)
			assert.Equal(t, createReq.Title, chat.Title)
			assert.Nil(t, chat.Avatar)
			return expectedChatInfo, nil
		})

	router := mux.NewRouter()
	delivery.NewChatController(router, mockChatUC, mockSessionClient)

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
	mockSessionClient := mocks.NewMockSessionServiceClient(ctrl)

	userID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	chatID := uuid.New()

	expectedChatInfo := &model.ChatInfo{
		ID:    chatID,
		Title: "Chat Info",
	}

	// Настраиваем ожидания для gRPC клиента
	mockSessionClient.EXPECT().
		CheckLogin(gomock.Any(), &authpb.CheckLoginRequest{SessionId: "00000000-0000-0000-0000-000000000001"}, gomock.Any()).
		Return(&authpb.CheckLoginResponse{UserId: userID.String()}, nil)

	mockChatUC.EXPECT().
		GetChatInfo(gomock.Any(), userID, chatID).
		Return(expectedChatInfo, nil)

	router := mux.NewRouter()
	delivery.NewChatController(router, mockChatUC, mockSessionClient)

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
