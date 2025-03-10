package utils

import (
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
)

func GetSessionCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func SetSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    sessionID,
		Path:     "/",
		Expires:  time.Now().Add(config.CookieDuration), // Кука на 3 часа
		HttpOnly: true,                                  // Защита от XSS
		SameSite: http.SameSiteNoneMode,
	})
}
