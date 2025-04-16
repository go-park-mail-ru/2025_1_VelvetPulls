package utils

import (
	"encoding/json"
	"errors"
	"net/http"
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
func SendJSONResponse(w http.ResponseWriter, statusCode int, v interface{}, success bool) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := JSONResponse{
		Status: success,
	}

	if success {
		// Если структура умеет себя санитизировать — даём ей это сделать
		if s, ok := v.(Sanitizable); ok {
			s.Sanitize()
		}

		// Если строка — санитизируем через строгий
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
		return errors.New("failed to encode JSON response: " + err.Error())
	}
	return nil
}
