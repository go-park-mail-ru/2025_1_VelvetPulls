package utils

import (
	"net/http"
	"time"
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
		Expires:  time.Now().Add(3 * time.Hour), // Кука на 3 часа
		HttpOnly: true,                          // Защита от XSS
		SameSite: http.SameSiteNoneMode,
	})
}
