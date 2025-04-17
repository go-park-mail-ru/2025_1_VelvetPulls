package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mock_usecase "github.com/go-park-mail-ru/2025_1_VelvetPulls/tests/delivery/mock"
)

func TestAuthMiddleware_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUC := mock_usecase.NewMockISessionUsecase(ctrl)
	userID := uuid.New()
	token := "valid-token"

	// Ожидаем вызов CheckLogin и возвращаем валидный UUID в виде строки.
	mockSessionUC.EXPECT().
		CheckLogin(gomock.Any(), token).
		Return(userID.String(), nil)

	// Формируем HTTP-запрос с кукой "token"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	rr := httptest.NewRecorder()

	called := false
	handler := middleware.AuthMiddleware(mockSessionUC)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		// Забираем userID из контекста
		ctxUserID := r.Context().Value(utils.USER_ID_KEY)
		assert.Equal(t, userID, ctxUserID)
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_NoCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUC := mock_usecase.NewMockISessionUsecase(ctrl)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler := middleware.AuthMiddleware(mockSessionUC)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler не должен быть вызван, если кука отсутствует")
	}))

	handler.ServeHTTP(rr, req)

	// Если кука отсутствует, должно возвращаться Unauthorized
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_InvalidSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUC := mock_usecase.NewMockISessionUsecase(ctrl)
	token := "invalid-token"

	// Ожидаем, что CheckLogin вернёт ошибку
	mockSessionUC.EXPECT().
		CheckLogin(gomock.Any(), token).
		Return("", errors.New("session not found"))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	rr := httptest.NewRecorder()

	handler := middleware.AuthMiddleware(mockSessionUC)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler не должен быть вызван, если сессия невалидна")
	}))

	handler.ServeHTTP(rr, req)

	// При ошибке сессии возвращаем BadRequest (как в реализации)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAuthMiddleware_InvalidUUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUC := mock_usecase.NewMockISessionUsecase(ctrl)
	token := "valid-token"

	// Возвращаем строку, которая не является корректным UUID
	mockSessionUC.EXPECT().
		CheckLogin(gomock.Any(), token).
		Return("not-a-uuid", nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	rr := httptest.NewRecorder()

	handler := middleware.AuthMiddleware(mockSessionUC)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler не должен быть вызван, если userID не валиден")
	}))

	handler.ServeHTTP(rr, req)

	// При ошибке преобразования UUID возвращаем BadRequest
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAuthMiddlewareWS_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUC := mock_usecase.NewMockISessionUsecase(ctrl)
	userID := uuid.New()
	token := "valid-token"

	mockSessionUC.EXPECT().
		CheckLogin(gomock.Any(), token).
		Return(userID.String(), nil)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	rr := httptest.NewRecorder()

	called := false
	handler := middleware.AuthMiddlewareWS(mockSessionUC)(func(w http.ResponseWriter, r *http.Request) {
		called = true
		id := r.Context().Value(utils.USER_ID_KEY)
		assert.Equal(t, userID, id)
		w.WriteHeader(http.StatusOK)
	})

	handler.ServeHTTP(rr, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddlewareWS_InvalidSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSessionUC := mock_usecase.NewMockISessionUsecase(ctrl)
	token := "invalid-token"

	mockSessionUC.EXPECT().
		CheckLogin(gomock.Any(), token).
		Return("", errors.New("session error"))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})
	rr := httptest.NewRecorder()

	handler := middleware.AuthMiddlewareWS(mockSessionUC)(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler не должен быть вызван при невалидной сессии")
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
