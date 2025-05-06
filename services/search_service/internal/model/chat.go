package model

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type Chat struct {
	ID           uuid.UUID    `json:"id" valid:"uuid"`
	AvatarPath   *string      `json:"avatar_path,omitempty"`
	Type         string       `json:"type" valid:"in(dialog|group|channel),required"`
	Title        string       `json:"title" valid:"required~Title is required,length(1|100)"`
	CreatedAt    string       `json:"created_at"`
	UpdatedAt    string       `json:"updated_at"`
	Participants []UserInChat `json:"participants"`
	LastMessage  *LastMessage `json:"last_message,omitempty"`
}

type UserInChat struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	AvatarPath *string   `json:"avatar_path"`
}

type RequestChat struct {
	Title string `json:"title" valid:"required~Title is required,length(1|100)"`
}

func (c *Chat) Validate() error {
	if _, err := govalidator.ValidateStruct(c); err != nil {
		return errors.Join(ErrValidation, errors.New("invalid chat data: "+err.Error()))
	}
	return nil
}

func (c *Chat) Sanitize() {
	c.Title = utils.SanitizeString(c.Title)
}

func (c *RequestChat) Sanitize() {
	c.Title = utils.SanitizeString(c.Title)
}
