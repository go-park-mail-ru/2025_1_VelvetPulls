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
	authpb "github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/proto"
	mocks "github.com/go-park-mail-ru/2025_1_VelvetPulls/tests/delivery/mock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := mocks.NewMockAuthServiceClient(ctrl)
	mockSessionClient := mocks.NewMockSessionServiceClient(ctrl)

	registerData := model.RegisterCredentials{
		Username:        "testuser123",
		Password:        "Password123!",
		ConfirmPassword: "Password123!",
		Phone:           "1234567890",
	}
	sessionID := "test-session-id"

	// Настройка gRPC моков
	mockAuthClient.EXPECT().
		RegisterUser(
			gomock.Any(),
			&authpb.RegisterUserRequest{
				Username:        registerData.Username,
				Password:        registerData.Password,
				ConfirmPassword: registerData.ConfirmPassword,
				Phone:           registerData.Phone,
			},
		).
		Return(&authpb.RegisterUserResponse{
			SessionId: sessionID,
		}, nil)

	router := mux.NewRouter()
	delivery.NewAuthController(router, mockAuthClient, mockSessionClient)

	body, err := json.Marshal(registerData)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp utils.JSONResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Status)
	assert.Equal(t, "Registration successful", resp.Data)

	cookies := rr.Result().Cookies()
	assert.NotEmpty(t, cookies)
	assert.Equal(t, "token", cookies[0].Name)
	assert.Equal(t, sessionID, cookies[0].Value)
}

func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := mocks.NewMockAuthServiceClient(ctrl)
	mockSessionClient := mocks.NewMockSessionServiceClient(ctrl)

	loginData := model.LoginCredentials{
		Username: "testuser123",
		Password: "Password123!",
	}
	sessionID := "test-session-id"

	// Настройка gRPC моков
	mockAuthClient.EXPECT().
		LoginUser(
			gomock.Any(),
			&authpb.LoginUserRequest{
				Username: loginData.Username,
				Password: loginData.Password,
			},
		).
		Return(&authpb.LoginUserResponse{
			SessionId: sessionID,
		}, nil)

	router := mux.NewRouter()
	delivery.NewAuthController(router, mockAuthClient, mockSessionClient)

	body, err := json.Marshal(loginData)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp utils.JSONResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Status)
	assert.Equal(t, "Login successful", resp.Data)

	cookies := rr.Result().Cookies()
	assert.NotEmpty(t, cookies)
	assert.Equal(t, "token", cookies[0].Name)
	assert.Equal(t, sessionID, cookies[0].Value)
}

func TestLogout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthClient := mocks.NewMockAuthServiceClient(ctrl)
	mockSessionClient := mocks.NewMockSessionServiceClient(ctrl)

	sessionID := "test-session-id"

	// Настройка gRPC моков
	mockAuthClient.EXPECT().
		LogoutUser(
			gomock.Any(),
			&authpb.LogoutUserRequest{
				SessionId: sessionID,
			},
		).
		Return(&emptypb.Empty{}, nil)

	router := mux.NewRouter()
	delivery.NewAuthController(router, mockAuthClient, mockSessionClient)

	req := httptest.NewRequest(http.MethodDelete, "/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: sessionID,
	})
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp utils.JSONResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Status)
	assert.Equal(t, "Logout successful", resp.Data)

	// Проверка удаления cookie
	cookies := rr.Result().Cookies()
	assert.NotEmpty(t, cookies)
	assert.Equal(t, "token", cookies[0].Name)
	assert.Equal(t, "", cookies[0].Value)
	assert.True(t, cookies[0].Expires.Before(time.Now()))
}
