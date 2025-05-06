package model

import "errors"

var (
	ErrValidation  = errors.New("validation error")
	ErrInvalidUUID = errors.New("invalid user UUID")
)
