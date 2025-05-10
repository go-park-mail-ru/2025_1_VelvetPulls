package model

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type ChatGroups struct {
	Dialogs  []Chat `json:"dialogs"`
	Groups   []Chat `json:"groups"`
	Channels []Chat `json:"channels"`
}

type Chat struct {
	ID          uuid.UUID    `json:"id" valid:"uuid"`
	AvatarPath  *string      `json:"avatar_path,omitempty"`
	Type        string       `json:"type" valid:"in(dialog|group|channel),required"`
	Title       string       `json:"title" valid:"required~Title is required,length(1|100)"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
	LastMessage *LastMessage `json:"last_message,omitempty"`
}

type RequestChat struct {
	Title string `json:"title" valid:"required~Title is required,length(1|100)"`
}

type UserInChat struct {
	ID         uuid.UUID `json:"id" valid:"uuid"`
	Username   string    `json:"username,omitempty" valid:"required,length(3|50)"`
	Name       *string   `json:"name,omitempty" valid:"length(0|100)"`
	AvatarPath *string   `json:"avatar_path,omitempty"`
	Role       *string   `json:"role" valid:"length(0|20)"`
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
