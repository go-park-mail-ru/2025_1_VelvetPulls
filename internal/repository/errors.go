package repository

import "errors"

var (
	ErrSessionNotFound      = errors.New("session not found")
	ErrInvalidInput         = errors.New("invalid input")
	ErrUserNotFound         = errors.New("user not found")
	ErrSelfContact          = errors.New("cannot add yourself as a contact")
	ErrRecordAlreadyExists  = errors.New("record already exists")
	ErrUpdateFailed         = errors.New("update failed")
	ErrChatNotFound         = errors.New("chat not found")
	ErrInvalidUUID          = errors.New("invalid UUID format")
	ErrEmptyField           = errors.New("empty required field")
	ErrDatabaseOperation    = errors.New("database operation failed")
	ErrDatabaseScan         = errors.New("failed to scan database row")
	ErrEmptyMessage         = errors.New("message body is empty")
	ErrSetNotifications     = errors.New("failed to update send_notifications status")
	ErrGetNotifications     = errors.New("failed to get send_notifications status")
	ErrContactAlreadyExists = errors.New("contact already exists")
)
