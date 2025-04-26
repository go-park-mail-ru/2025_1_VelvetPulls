package repository

import "errors"

var (
	ErrEmptyField          = errors.New("empty required field")
	ErrInvalidInput        = errors.New("invalid input")
	ErrRecordAlreadyExists = errors.New("record already exists")
	ErrNotFound            = errors.New("record not found")
	ErrDatabaseOperation   = errors.New("database operation failed")
	ErrInvalidUUID         = errors.New("invalid UUID format")
	ErrInvalidRating       = errors.New("invalid rating value")
	ErrNoStatistics        = errors.New("no statistics available")
	ErrUserNotActive       = errors.New("user has no activity records")
)
