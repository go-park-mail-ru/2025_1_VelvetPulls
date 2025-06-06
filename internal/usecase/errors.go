package usecase

import "errors"

var (
	ErrUsernameIsTaken = errors.New("username is already taken")
	ErrPhoneIsTaken    = errors.New("phone number is already taken")
	ErrHashPassword    = errors.New("failed to hash password")
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidPassword = errors.New("invalid password")

	ErrPermissionDenied        = errors.New("permission denied")
	ErrDialogUpdateForbidden   = errors.New("cannot update a dialog")
	ErrOnlyOwnerCanModify      = errors.New("only chat owner can modify chat")
	ErrDialogAddUsers          = errors.New("cannot add users to a dialog")
	ErrDialogDeleteUsers       = errors.New("cannot delete users from a dialog")
	ErrChatCreationFailed      = errors.New("failed to create chat")
	ErrAddOwnerToDialog        = errors.New("failed to add owner to dialog")
	ErrAddParticipantToDialog  = errors.New("failed to add participant to dialog as owner")
	ErrAddOwnerToGroup         = errors.New("failed to add owner to group")
	ErrOnlyOwnerCanDelete      = errors.New("only chat owner can delete chat")
	ErrOnlyOwnerCanAddUsers    = errors.New("only chat owner can add users")
	ErrNotChannel              = errors.New("chat type is not channel")
	ErrOnlyOwnerCanDeleteUsers = errors.New("only chat owner can delete users")

	ErrMessageValidationFailed = errors.New("message validation failed")
	ErrMessageCreationFailed   = errors.New("failed to create message")
	ErrMessageNotFound         = errors.New("message not found")
	ErrMessageAccessDenied     = errors.New("user is not the author of the message")
	ErrMessageUpdateFailed     = errors.New("failed to update message")
	ErrMessageDeleteFailed     = errors.New("failed to delete message")

	ErrMessagePublishFailed = errors.New("failed to publish message event")
	ErrChatPublishFailed    = errors.New("failed to publish chat event")
)
