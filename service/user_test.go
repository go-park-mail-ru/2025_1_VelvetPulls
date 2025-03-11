package service

import (
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/stretchr/testify/require"
)

func TestRegisterUser(t *testing.T) {
	negativeTests := []struct {
		testName    string
		credentials model.RegisterCredentials
		response    string
		err         error
	}{
		{
			testName: "Different passwords",
			credentials: model.RegisterCredentials{
				Password:        "password",
				ConfirmPassword: "my_password",
			},
			response: "",
			err:      apperrors.ErrPasswordsDoNotMatch,
		},
		{
			testName: "Invalid password",
			credentials: model.RegisterCredentials{
				Password:        "pass word",
				ConfirmPassword: "pass word",
			},
			response: "",
			err:      apperrors.ErrInvalidPassword,
		},
		{
			testName: "Invalid phone",
			credentials: model.RegisterCredentials{
				Password:        "password",
				ConfirmPassword: "password",
				Phone:           "abc",
			},
			response: "",
			err:      apperrors.ErrInvalidPhoneFormat,
		},
		{
			testName: "Invalid username",
			credentials: model.RegisterCredentials{
				Username:        "Максимка",
				Password:        "password",
				ConfirmPassword: "password",
				Phone:           "+79251385523",
			},
			response: "",
			err:      apperrors.ErrInvalidUsername,
		},
		{
			testName: "User already exists",
			credentials: model.RegisterCredentials{
				Username:        "ruslantus228",
				Password:        "password",
				ConfirmPassword: "password",
				Phone:           "+79251385523",
			},
			response: "",
			err:      apperrors.ErrUserAlreadyExists,
		},
		{
			testName: "Phone number is already taken",
			credentials: model.RegisterCredentials{
				Username:        "lolkekcheburek",
				Password:        "password",
				ConfirmPassword: "password",
				Phone:           "+79128234765",
			},
			response: "",
			err:      apperrors.ErrPhoneTaken,
		},
		{
			testName: "Phone number is already taken",
			credentials: model.RegisterCredentials{
				Username:        "lolkekcheburek",
				Password:        "password",
				ConfirmPassword: "password",
				Phone:           "+79128234765",
			},
			response: "",
			err:      apperrors.ErrPhoneTaken,
		},
	}

	for _, tt := range negativeTests {
		t.Run(tt.testName, func(t *testing.T) {
			response, err := RegisterUser(tt.credentials)

			require.Equal(t, tt.response, response)
			require.Equal(t, tt.err, err)
		})
	}

	positiveTests := []struct {
		testName    string
		credentials model.RegisterCredentials
		response    string
		err         error
	}{
		{
			testName: "Successful registration",
			credentials: model.RegisterCredentials{
				Username:        "lolkekcheburek",
				Password:        "password",
				ConfirmPassword: "password",
				Phone:           "+79123457235",
			},
			response: "",
			err:      nil,
		},
	}

	for _, tt := range positiveTests {
		t.Run(tt.testName, func(t *testing.T) {
			response, err := RegisterUser(tt.credentials)

			require.NotEqual(t, tt.response, response)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestLoginUser(t *testing.T) {
	negativeTests := []struct {
		testName    string
		credentials model.LoginCredentials
		response    string
		err         error
	}{
		{
			testName: "User is not found",
			credentials: model.LoginCredentials{
				Username: "kefteme",
			},
			response: "",
			err:      apperrors.ErrUserNotFound,
		},
		{
			testName: "User is not found",
			credentials: model.LoginCredentials{
				Username: "ruslantus228",
				Password: "not_my_password",
			},
			response: "",
			err:      apperrors.ErrInvalidCredentials,
		},
	}

	for _, tt := range negativeTests {
		t.Run(tt.testName, func(t *testing.T) {
			response, err := LoginUser(tt.credentials)

			require.Equal(t, tt.response, response)
			require.Equal(t, tt.err, err)
		})
	}

	positiveTests := []struct {
		testName    string
		credentials model.LoginCredentials
		response    string
		err         error
	}{
		{
			testName: "Successful authorization",
			credentials: model.LoginCredentials{
				Username: "ruslantus228",
				Password: "qwerty",
			},
			response: "",
			err:      nil,
		},
	}

	for _, tt := range positiveTests {
		t.Run(tt.testName, func(t *testing.T) {
			response, err := LoginUser(tt.credentials)

			require.NotEqual(t, tt.response, response)
			require.Equal(t, tt.err, err)
		})
	}
}
