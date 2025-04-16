package utils_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestParseJSONRequest_Success(t *testing.T) {
	body := `{"name":"John","age":30}`
	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewBufferString(body)))

	var data testStruct
	err := utils.ParseJSONRequest(req, &data)

	assert.NoError(t, err)
	assert.Equal(t, "John", data.Name)
	assert.Equal(t, 30, data.Age)
}

func TestParseJSONRequest_InvalidJSON(t *testing.T) {
	body := `{"name": "Alice",` // некорректный JSON
	req := httptest.NewRequest(http.MethodPost, "/", io.NopCloser(bytes.NewBufferString(body)))

	var data testStruct
	err := utils.ParseJSONRequest(req, &data)

	assert.Error(t, err)
}

func TestSendJSONResponse_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	data := testStruct{Name: "Bob", Age: 25}

	err := utils.SendJSONResponse(rr, http.StatusOK, data, true)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)

	var response utils.JSONResponse
	json.NewDecoder(rr.Body).Decode(&response)

	assert.True(t, response.Status)
	assert.Equal(t, float64(25), response.Data.(map[string]interface{})["age"])
}

func TestSendJSONResponse_ErrorFromError(t *testing.T) {
	rr := httptest.NewRecorder()
	err := errors.New("something went wrong")

	utils.SendJSONResponse(rr, http.StatusInternalServerError, err, false)

	var response utils.JSONResponse
	json.NewDecoder(rr.Body).Decode(&response)

	assert.False(t, response.Status)
	assert.Equal(t, "something went wrong", response.Error)
}

func TestSendJSONResponse_ErrorFromString(t *testing.T) {
	rr := httptest.NewRecorder()
	utils.SendJSONResponse(rr, http.StatusBadRequest, "bad request", false)

	var response utils.JSONResponse
	json.NewDecoder(rr.Body).Decode(&response)

	assert.False(t, response.Status)
	assert.Equal(t, "bad request", response.Error)
}

func TestSendJSONResponse_UnknownError(t *testing.T) {
	rr := httptest.NewRecorder()
	utils.SendJSONResponse(rr, http.StatusBadRequest, 123, false)

	var response utils.JSONResponse
	json.NewDecoder(rr.Body).Decode(&response)

	assert.False(t, response.Status)
	assert.Equal(t, "unknown error", response.Error)
}
