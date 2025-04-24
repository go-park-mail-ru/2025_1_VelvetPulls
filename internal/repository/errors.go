package repository

import "errors"

var (
	ErrSessionNotFound     = errors.New("session not found")
	ErrInvalidInput        = errors.New("invalid input")
	ErrUserNotFound        = errors.New("user not found")
	ErrSelfContact         = errors.New("cannot add yourself as a contact")
	ErrRecordAlreadyExists = errors.New("record already exists")
	ErrUpdateFailed        = errors.New("update failed")
	ErrChatNotFound        = errors.New("chat not found")
	ErrInvalidUUID         = errors.New("invalid UUID format")
	ErrEmptyField          = errors.New("empty required field")
	ErrDatabaseOperation   = errors.New("database operation failed")
	ErrDatabaseScan        = errors.New("failed to scan database row")
	ErrEmptyMessage        = errors.New("message body is empty")
	ErrFileNotFound        = errors.New("file not found")
	ErrFileAccessDenied    = errors.New("access to file denied")
	ErrFileStorageFailed   = errors.New("file storage operation failed")
)
