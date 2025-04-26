package usecase

import "errors"

var (
	ErrInvalidInput        = errors.New("invalid input data")
	ErrNotFound            = errors.New("not found")
	ErrDatabaseOperation   = errors.New("database operation error")
	ErrInternalServerError = errors.New("internal server error")
)
