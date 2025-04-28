package utils

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// JSONResponse - структура для унифицированного ответа.
type JSONResponse struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// ParseJSONRequest парсит JSON из тела запроса в переданную структуру.
func ParseJSONRequest(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(v)
}

// SendJSONResponse отправляет JSON-ответ с полем `status`.
func SendJSONResponse(w http.ResponseWriter, r *http.Request, statusCode int, v interface{}, success bool) {
	logger := GetLoggerFromCtx(r.Context())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := JSONResponse{
		Status: success,
	}

	if success {
		if s, ok := v.(Sanitizable); ok {
			s.Sanitize()
		}

		if str, ok := v.(string); ok {
			response.Data = SanitizeString(str)
		} else {
			response.Data = v
		}
	} else {
		if err, ok := v.(error); ok {
			response.Error = SanitizeString(err.Error())
		} else if str, ok := v.(string); ok {
			response.Error = SanitizeString(str)
		} else {
			response.Error = "unknown error"
		}
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		if logger != nil {
			logger.Error("failed to encode JSON response", zap.Error(err))
		}
	}
}
