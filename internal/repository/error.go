package repository

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrAvatarNotFound = errors.New("avatar not found")
)
