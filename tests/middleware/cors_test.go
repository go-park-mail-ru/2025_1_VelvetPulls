package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCorsMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "Regular request",
			method:         "GET",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "OPTIONS request",
			method:         "OPTIONS",
			expectedStatus: http.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest(tt.method, "/", nil)
			rr := httptest.NewRecorder()

			middleware.CorsMiddleware(handler).ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, config.Cors.AllowedOrigin, rr.Header().Get("Access-Control-Allow-Origin"))
			assert.Equal(t, "true", rr.Header().Get("Access-Control-Allow-Credentials"))
			assert.Equal(t, config.Cors.AllowedMethods, rr.Header().Get("Access-Control-Allow-Methods"))

			if tt.method == "OPTIONS" {
				assert.Contains(t, rr.Header().Get("Access-Control-Allow-Headers"), "X-CSRF-Token")
				assert.Equal(t, "X-CSRF-Token", rr.Header().Get("Access-Control-Expose-Headers"))
			}
		})
	}
}
