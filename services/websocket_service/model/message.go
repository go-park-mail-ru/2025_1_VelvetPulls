package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID              uuid.UUID  `json:"id,omitempty"`
	ParentMessageID *uuid.UUID `json:"parent_message_id,omitempty"`
	ChatID          uuid.UUID  `json:"chat_id,omitempty"`
	UserID          uuid.UUID  `json:"user_id,omitempty"`
	Body            string     `json:"body,omitempty"`
	SentAt          time.Time  `json:"sent_at,omitempty"`
	IsRedacted      bool       `json:"is_redacted,omitempty"`
	AvatarPath      *string    `json:"avatar_path,omitempty"`
	Username        string     `json:"user,omitempty"`
}

type ChatInfo struct {
	ID         uuid.UUID    `json:"id" valid:"uuid"`
	AvatarPath *string      `json:"avatar_path,omitempty"`
	Type       string       `json:"type" valid:"in(dialog|group|channel)"`
	Title      string       `json:"title" valid:"length(1|100)"`
	CountUsers int          `json:"count_users" valid:"range(0|5000)"`
	Users      []UserInChat `json:"users"`
}

type UserInChat struct {
	ID         uuid.UUID `json:"id" valid:"uuid"`
	Username   string    `json:"username,omitempty" valid:"required,length(3|50)"`
	Name       *string   `json:"name,omitempty" valid:"length(0|100)"`
	AvatarPath *string   `json:"avatar_path,omitempty"`
	Role       *string   `json:"role" valid:"length(0|20)"`
}
