package repository

import (
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/stretchr/testify/require"
)

func TestGetUserByUsername(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:          "User found: ruslantus228",
			username:      "ruslantus228",
			expectedUser:  users[0],
			expectedError: nil,
		},
		{
			name:          "User not found: charlie",
			username:      "charlie",
			expectedUser:  nil,
			expectedError: apperrors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := GetUserByUsername(tt.username)

			require.Equal(t, tt.expectedUser, user)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:          "User found: rumail@mail.ru",
			email:         "rumail@mail.ru",
			expectedUser:  users[0],
			expectedError: nil,
		},
		{
			name:          "User not found: notexisted@mail.ru",
			email:         "notexisted@mail.ru",
			expectedUser:  nil,
			expectedError: apperrors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := GetUserByEmail(tt.email)

			require.Equal(t, tt.expectedUser, user)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestGetUserByPhone(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		expectedUser  *model.User
		expectedError error
	}{
		{
			name:          "User found: +79128234765",
			email:         "+79128234765",
			expectedUser:  users[0],
			expectedError: nil,
		},
		{
			name:          "User not found: 8-800-555-35-35",
			email:         "8-800-555-35-35",
			expectedUser:  nil,
			expectedError: apperrors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := GetUserByPhone(tt.email)

			require.Equal(t, tt.expectedUser, user)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name          string
		user          *model.User
		expectedError error
	}{
		{
			name: "New user",
			user: &model.User{
				Username: "Dyadya_Vo1odya",
			},
			expectedError: nil,
		},
		{
			name: "Username is already taken",
			user: &model.User{
				Username: "ruslantus228",
			},
			expectedError: apperrors.ErrUsernameTaken,
		},
		{
			name: "Phone number is already taken",
			user: &model.User{
				Phone: "+79128234765",
			},
			expectedError: apperrors.ErrPhoneTaken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CreateUser(tt.user)

			require.Equal(t, tt.expectedError, err)
		})
	}
}
