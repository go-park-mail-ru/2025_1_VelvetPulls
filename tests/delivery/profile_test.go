package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestGetSelfProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUC := mocks.NewMockIUserUsecase(ctrl)
	mockSessionClient := mocks.NewMockSessionServiceClient(ctrl)

	userID := uuid.New()
	expectedProfile := &model.GetUserProfile{
		Username: "testuser",
		Name:     "1234567890",
	}

	// Настройка gRPC моков
	mockSessionClient.EXPECT().
		CheckLogin(
			gomock.Any(),
			&authpb.CheckLoginRequest{SessionId: "valid-token"},
			gomock.Any(),
		).
		Return(&authpb.CheckLoginResponse{UserId: userID.String()}, nil)

	mockUserUC.EXPECT().
		GetUserProfileByID(gomock.Any(), userID).
		Return(expectedProfile, nil)

	router := mux.NewRouter()
	delivery.NewUserController(router, mockUserUC, mockSessionClient)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
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

	var profile model.GetUserProfile
	jsonData, err := json.Marshal(resp.Data)
	assert.NoError(t, err)
	err = json.Unmarshal(jsonData, &profile)
	assert.NoError(t, err)

	assert.Equal(t, *expectedProfile, profile)
}

func TestGetProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUC := mocks.NewMockIUserUsecase(ctrl)
	mockSessionClient := mocks.NewMockSessionServiceClient(ctrl)

	username := "testuser"
	expectedProfile := &model.GetUserProfile{
		Username: username,
		Name:     "1234567890",
	}

	// Настройка gRPC моков
	mockSessionClient.EXPECT().
		CheckLogin(
			gomock.Any(),
			&authpb.CheckLoginRequest{SessionId: "valid-token"},
			gomock.Any(),
		).
		Return(&authpb.CheckLoginResponse{UserId: uuid.New().String()}, nil)

	mockUserUC.EXPECT().
		GetUserProfileByUsername(gomock.Any(), username).
		Return(expectedProfile, nil)

	router := mux.NewRouter()
	delivery.NewUserController(router, mockUserUC, mockSessionClient)

	req := httptest.NewRequest(http.MethodGet, "/profile/"+username, nil)
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

	var profile model.GetUserProfile
	jsonData, err := json.Marshal(resp.Data)
	assert.NoError(t, err)
	err = json.Unmarshal(jsonData, &profile)
	assert.NoError(t, err)

	assert.Equal(t, *expectedProfile, profile)
}

func TestUpdateSelfProfile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUC := mocks.NewMockIUserUsecase(ctrl)
	mockSessionClient := mocks.NewMockSessionServiceClient(ctrl)

	userID := uuid.New()
	Name := "NewFirstName"
	updateData := model.UpdateUserProfile{
		ID:   userID,
		Name: &Name,
	}

	// Настройка gRPC моков
	mockSessionClient.EXPECT().
		CheckLogin(
			gomock.Any(),
			&authpb.CheckLoginRequest{SessionId: "valid-token"},
			gomock.Any(),
		).
		Return(&authpb.CheckLoginResponse{UserId: userID.String()}, nil)

	mockUserUC.EXPECT().
		UpdateUserProfile(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, profile *model.UpdateUserProfile) error {
			assert.Equal(t, userID, profile.ID)
			assert.Equal(t, *updateData.Name, *profile.Name)
			return nil
		})

	router := mux.NewRouter()
	delivery.NewUserController(router, mockUserUC, mockSessionClient)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	profileData, err := json.Marshal(updateData)
	assert.NoError(t, err)

	writer.WriteField("profile_data", string(profileData))
	writer.Close()

	req := httptest.NewRequest(http.MethodPut, "/profile", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "valid-token",
	})
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp utils.JSONResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Status)
	assert.Equal(t, "Profile updated successfully", resp.Data)
}

func TestUpdateSelfProfile_InvalidData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUC := mocks.NewMockIUserUsecase(ctrl)
	mockSessionClient := mocks.NewMockSessionServiceClient(ctrl)

	// Настройка gRPC моков
	mockSessionClient.EXPECT().
		CheckLogin(
			gomock.Any(),
			&authpb.CheckLoginRequest{SessionId: "valid-token"},
			gomock.Any(),
		).
		Return(&authpb.CheckLoginResponse{UserId: uuid.New().String()}, nil)

	router := mux.NewRouter()
	delivery.NewUserController(router, mockUserUC, mockSessionClient)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("profile_data", "invalid json")
	writer.Close()

	req := httptest.NewRequest(http.MethodPut, "/profile", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "valid-token",
	})
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetProfile_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUC := mocks.NewMockIUserUsecase(ctrl)
	mockSessionClient := mocks.NewMockSessionServiceClient(ctrl)

	router := mux.NewRouter()
	delivery.NewUserController(router, mockUserUC, mockSessionClient)

	req := httptest.NewRequest(http.MethodGet, "/profile/testuser", nil)
	req = req.WithContext(context.WithValue(req.Context(), utils.LOGGER_ID_KEY, zap.NewNop()))

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
