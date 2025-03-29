package model

import (
	"time"

	"github.com/google/uuid"
)

// TODO: Переделать под новую структуру бд

type ChatType string

const (
	ChatTypeDialog  ChatType = "dialog"
	ChatTypeGroup   ChatType = "group"
	ChatTypeChannel ChatType = "channel"
)

type Chat struct {
	ID         uuid.UUID `json:"id"`
	AvatarPath *string   `json:"avatar_path,omitempty"`
	Type       ChatType  `json:"type"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
