package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
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

type LastMessage struct {
	ID       uuid.UUID `json:"id,omitempty"`
	UserID   uuid.UUID `json:"user_id,omitempty"`
	Body     string    `json:"body,omitempty"`
	SentAt   time.Time `json:"sent_at,omitempty"`
	Username string    `json:"user,omitempty"`
}

type MessageInput struct {
	Message string `json:"message" valid:"required,length(1|1000)"`
}

func (m *MessageInput) Validate() error {
	if _, err := govalidator.ValidateStruct(m); err != nil {
		return errors.Join(ErrValidation, fmt.Errorf("invalid message input: %w", err))
	}
	return nil
}

func (m *Message) Sanitize() {
	m.Body = utils.SanitizeString(m.Body)
	m.Username = utils.SanitizeString(m.Username)
}

func (mi *MessageInput) Sanitize() {
	mi.Message = utils.SanitizeString(mi.Message)
}
