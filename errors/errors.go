package errors

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailTaken        = errors.New("email already registered")
	ErrPhoneTaken        = errors.New("phone number already registered")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserCreation      = errors.New("user creation error")
	ErrInvalidParams     = errors.New("invalid params")
)
