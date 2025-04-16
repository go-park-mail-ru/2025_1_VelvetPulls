package middleware

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", config.Cors.AllowedOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", config.Cors.AllowedMethods)
		w.Header().Set("Access-Control-Allow-Headers", config.Cors.AllowedHeaders)

		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Headers",
				config.Cors.AllowedHeaders+", X-CSRF-Token")
			w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Access-Control-Expose-Headers", "X-CSRF-Token")

		next.ServeHTTP(w, r)
	})
}
