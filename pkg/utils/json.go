package utils

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// JSONResponse - структура для унифицированного ответа.
type JSONResponse struct {
	Status bool            `json:"status"`
	Data   json.RawMessage `json:"data"`
	Error  string          `json:"error,omitempty"`
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
		switch data := v.(type) {
		case []byte:
			// Если данные уже в формате JSON, используем их как есть
			response.Data = json.RawMessage(data)

		case string:
			// Санируем и преобразуем строку
			sanitized := SanitizeString(data)
			response.Data = json.RawMessage(`"` + sanitized + `"`)

		default:
			// Санизация и стандартное кодирование
			if s, ok := v.(Sanitizable); ok {
				s.Sanitize()
			}

			// Кодируем объект в JSON
			jsonData, err := json.Marshal(v)
			if err != nil {
				if logger != nil {
					logger.Error("failed to marshal data", zap.Error(err))
				}
				jsonData = []byte("null")
			}
			response.Data = jsonData
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

	// Кодируем финальный ответ
	if err := json.NewEncoder(w).Encode(response); err != nil {
		if logger != nil {
			logger.Error("failed to encode JSON response", zap.Error(err))
		}
	}
}
