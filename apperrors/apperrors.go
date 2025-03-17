package apperrors

import "errors"

// Ошибки, связанные с пользователями.
var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUsernameTaken       = errors.New("username already taken")
	ErrEmailTaken          = errors.New("email already registered")
	ErrPhoneTaken          = errors.New("phone number already registered")
	ErrUserCreation        = errors.New("user creation error")
	ErrInvalidCredentials  = errors.New("wrong password or username")
	ErrInvalidPhoneFormat  = errors.New("invalid phone format")
	ErrInvalidPassword     = errors.New("password must be between 8 and 32 characters long")
	ErrInvalidUsername     = errors.New("username must be between 3 and 20 characters long and can only contain letters, digits, and underscores")
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")
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

// Ошибки, связанные с сервером.
var (
	ErrInternalServer = errors.New("internal server error")
)
