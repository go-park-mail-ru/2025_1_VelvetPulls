package model

import (
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type Contact struct {
	ID        uuid.UUID `json:"id"`
	Name      *string   `json:"name,omitempty"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_path,omitempty"`
}

type RequestContact struct {
	Username string `json:"username"`
}

func (c *Contact) Sanitize() {
	c.Username = utils.SanitizeString(c.Username)
	if c.Name != nil {
		s := utils.SanitizeString(*c.Name)
		c.Name = &s
	}
	if c.AvatarURL != nil {
		s := utils.SanitizeString(*c.AvatarURL)
		c.AvatarURL = &s
	}
}

func (r *RequestContact) Sanitize() {
	r.Username = utils.SanitizeString(r.Username)
}
