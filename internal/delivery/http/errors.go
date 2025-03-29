package http

import "errors"

var (
	ErrInvalidRequestData = errors.New("invalid request data")
	ErrSessionNotFound    = errors.New("session not found")
)
