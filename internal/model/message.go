package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type MsgType string

const (
	NewMessage      MsgType = "message"
	AddUserInChat   MsgType = "addUserInChat"
	DelUserFromChat MsgType = "delUserFromChat"
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

type SendMessage struct {
	MessageType MsgType     `json:"messageType"`
	Payload     interface{} `json:"payload"`
}

func (sm *SendMessage) Validate() error {
	if sm.Payload == nil {
		return errors.Join(ErrValidation, errors.New("payload is required"))
	}

	switch sm.MessageType {
	case NewMessage, AddUserInChat, DelUserFromChat:
		// ok
	default:
		return errors.Join(ErrValidation, fmt.Errorf("unknown message type: %s", sm.MessageType))
	}

	if sm.MessageType == NewMessage {
		msg, ok := sm.Payload.(Message)
		if !ok {
			return errors.Join(ErrValidation, errors.New("invalid payload type for NewMessage"))
		}

		if _, err := govalidator.ValidateStruct(msg); err != nil {
			return errors.Join(ErrValidation, errors.New("invalid message payload: "+err.Error()))
		}
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

func (sm *SendMessage) Sanitize() {
	if sm.MessageType == NewMessage {
		if msg, ok := sm.Payload.(Message); ok {
			msg.Sanitize()
			sm.Payload = msg
		}
	}
}
