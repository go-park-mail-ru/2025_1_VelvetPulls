package utils_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.CookieDuration = 3 * time.Hour // чтобы не было nil во время тестов
}

func TestSetSessionCookie(t *testing.T) {
	rr := httptest.NewRecorder()
	sessionID := "test-session-id"

	utils.SetSessionCookie(rr, sessionID)

	result := rr.Result()
	cookies := result.Cookies()
	assert.Len(t, cookies, 1)

	c := cookies[0]
	assert.Equal(t, "token", c.Name)
	assert.Equal(t, sessionID, c.Value)
	assert.Equal(t, "/", c.Path)
	assert.True(t, c.HttpOnly)
	assert.WithinDuration(t, time.Now().Add(3*time.Hour), c.Expires, time.Second)
}

func TestGetSessionCookie(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: "some-session-id",
	})

	sessionID, err := utils.GetSessionCookie(req)
	assert.NoError(t, err)
	assert.Equal(t, "some-session-id", sessionID)
}

func TestGetSessionCookie_Missing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	sessionID, err := utils.GetSessionCookie(req)
	assert.Error(t, err)
	assert.Empty(t, sessionID)
}

func TestDeleteSessionCookie(t *testing.T) {
	rr := httptest.NewRecorder()

	utils.DeleteSessionCookie(rr)

	result := rr.Result()
	cookies := result.Cookies()
	assert.Len(t, cookies, 1)

	c := cookies[0]
	assert.Equal(t, "token", c.Name)
	assert.Equal(t, "", c.Value)
	assert.Equal(t, "/", c.Path)
	assert.True(t, c.HttpOnly)
	assert.True(t, c.Expires.Before(time.Now()))
}
