package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/middleware"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestAccessLogMiddleware(t *testing.T) {
	tests := []struct {
		name         string
		contentType  string
		responseBody string
		expectedLog  bool
	}{
		{
			name:         "JSON response",
			contentType:  "application/json",
			responseBody: `{"message":"test"}`,
			expectedLog:  true,
		},
		{
			name:         "Non-JSON response",
			contentType:  "text/plain",
			responseBody: "test",
			expectedLog:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tt.contentType)
				w.Write([]byte(tt.responseBody))
			})

			req := httptest.NewRequest("GET", "/", nil)
			req = req.WithContext(context.WithValue(req.Context(), utils.REQUEST_ID_KEY, "test-request-id"))

			rr := httptest.NewRecorder()
			middleware.AccessLogMiddleware(handler).ServeHTTP(rr, req)

			assert.Equal(t, http.StatusOK, rr.Code)
		})
	}
}
