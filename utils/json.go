package utils

import (
	"encoding/json"
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
		response.Data = v
	} else {
		if err, ok := v.(error); ok {
			response.Error = err.Error()
		} else if str, ok := v.(string); ok {
			response.Error = str
		} else {
			response.Error = "unknown error"
		}
	}

	return json.NewEncoder(w).Encode(response)
}
