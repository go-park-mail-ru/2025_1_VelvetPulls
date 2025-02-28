package errors

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUsernameTaken        = errors.New("no such user with this username")
	ErrEmailTaken           = errors.New("email already registered")
	ErrPhoneTaken           = errors.New("phone number already registered")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrSessionAlreadyExists = errors.New("session already exists")
	ErrUserCreation         = errors.New("user creation error")
	ErrInvalidParams        = errors.New("invalid params")
	ErrSessionNotFound      = errors.New("session not found")
	ErrInvalidCredentials   = errors.New("wrong password or username")
)
