package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/gorilla/csrf"
)

func CSRF(secure bool, authKey []byte) func(http.Handler) http.Handler {
	return csrf.Protect(
		authKey,
		csrf.Secure(secure), //false - http
		csrf.Path("/"),
		csrf.HttpOnly(true),
		csrf.FieldName("csrf_token"),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "CSRF token invalid", http.StatusForbidden)
		})),
		csrf.TrustedOrigins([]string{config.Cors.AllowedOrigin}),
	)
}
