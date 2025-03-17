package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name   string
		data   map[string]string
		status int
	}{
		{
			name: "Successful registration",
			data: map[string]string{
				"username":         "gojsfullstack",
				"confirm_password": "lolkekcheburek",
				"password":         "lolkekcheburek",
				"phone":            "+79991234567",
			},
			status: http.StatusCreated,
		},
		{
			name: "Invalid body",
			data: map[string]string{
				"field":         "gojsfullstack",
				"another_field": "lolkekcheburek",
			},
			status: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.data)
			require.Equal(t, nil, err)

			body := bytes.NewReader(data)

			r := httptest.NewRequest("POST", "/api/register/", body)
			w := httptest.NewRecorder()

			Register(w, r)

			require.Equal(t, test.status, w.Code)
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name   string
		data   map[string]string
		status int
	}{
		{
			name: "Successful authorization",
			data: map[string]string{
				"username": "ruslantus228",
				"password": "qwerty",
			},
			status: http.StatusOK,
		},
		{
			name: "Invalid body",
			data: map[string]string{
				"field":         "gojsfullstack",
				"another_field": "lolkekcheburek",
			},
			status: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.data)
			require.Equal(t, nil, err)

			body := bytes.NewReader(data)

			r := httptest.NewRequest("POST", "/api/login/", body)
			w := httptest.NewRecorder()

			Login(w, r)

			require.Equal(t, test.status, w.Code)
		})
	}
}
