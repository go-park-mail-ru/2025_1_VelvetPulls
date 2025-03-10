package apperrors

import "errors"

// Ошибки, связанные с пользователями.
var (
	ErrUserNotFound        = errors.New("user not found")
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")

	ErrUsernameTaken      = errors.New("username already taken")
	ErrEmailTaken         = errors.New("email already registered")
	ErrPhoneTaken         = errors.New("phone number already registered")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserCreation       = errors.New("user creation error")
	ErrInvalidCredentials = errors.New("wrong password or username")
	ErrInvalidPhoneFormat = errors.New("invalid phone format")
	ErrInvalidPassword    = errors.New("password must be at least 8 characters long")
	ErrInvalidUsername    = errors.New("username must be at least 3 characters long")
)

// Ошибки, связанные с сессиями.
var (
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionAlreadyExists = errors.New("session already exists")
)

// Ошибки, связанные с чатами.
var (
	ErrChatNotFound      = errors.New("chat not found")
	ErrChatAlreadyExists = errors.New("chat already exists")
	ErrUserNotInChat     = errors.New("user is not a member of this chat")
)

// Ошибки, связанные с параметрами.
var (
	ErrInvalidParams = errors.New("invalid parameters")
)
