package repository

import "errors"

var (
	ErrSearchMessages      = errors.New("search messages query failed")
	ErrSearchMessagesCount = errors.New("get messages count failed")
	ErrSearchUsers         = errors.New("search users query failed")
	ErrSearchContacts      = errors.New("search contacts query failed")
	ErrSearchUserChats     = errors.New("search user chats query failed")
)
