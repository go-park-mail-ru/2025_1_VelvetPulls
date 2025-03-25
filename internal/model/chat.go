package model

import (
	"time"
)

// TODO: Переделать под новую структуру бд

type ChatType string

const (
	ChatTypeDialog  ChatType = "dialog"
	ChatTypeGroup   ChatType = "group"
	ChatTypeChannel ChatType = "channel"
)

type Chat struct {
	ID            int64     `json:"id"`
	OwnerUsername string    `json:"owner_username"`
	Type          ChatType  `json:"type"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Members       []int64   `json:"members"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
