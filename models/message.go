package models

import (
	"time"
)

type Message struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	UserID    int64     `json:"user_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
