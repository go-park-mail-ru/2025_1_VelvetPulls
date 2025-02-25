package errors

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserCreation      = errors.New("user creation error")
	ErrInvalidParams     = errors.New("invalid params")
)
