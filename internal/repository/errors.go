package repository

import "errors"

var (
	ErrSessionNotFound     = errors.New("session not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrRecordAlreadyExists = errors.New("record already exists")
	ErrUpdateFailed        = errors.New("update failed")
	ErrInvalidUUID         = errors.New("invalid UUID format")
	ErrEmptyField          = errors.New("empty required field")
	ErrDatabaseOperation   = errors.New("database operation failed")
)
