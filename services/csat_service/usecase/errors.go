package usecase

import "errors"

var (
	ErrInvalidInput        = errors.New("invalid input data")
	ErrNotFound            = errors.New("not found")
	ErrInternalServerError = errors.New("internal server error")
)
