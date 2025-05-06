package usecase

import "errors"

var (
	ErrUsernameIsTaken = errors.New("username is already taken")
	ErrPhoneIsTaken    = errors.New("phone number is already taken")
	ErrHashPassword    = errors.New("failed to hash password")
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")
)
