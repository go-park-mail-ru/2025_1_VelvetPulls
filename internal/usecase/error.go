package usecase

import "errors"

var (
	ErrBadAvatarSize = errors.New("avatar size exceeds the limit of 2MB")
	ErrBadAvatarType = errors.New("unsupported avatar file type")
	ErrReadAvatar    = errors.New("can't read avatar")
)
