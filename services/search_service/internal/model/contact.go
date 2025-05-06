package model

import (
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type Contact struct {
	ID        uuid.UUID `json:"id"`
	FirstName *string   `json:"first_name,omitempty"`
	LastName  *string   `json:"last_name,omitempty"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_path,omitempty"`
}

type RequestContact struct {
	Username string `json:"username"`
}

func (c *Contact) Sanitize() {
	c.Username = utils.SanitizeString(c.Username)
	if c.FirstName != nil {
		s := utils.SanitizeString(*c.FirstName)
		c.FirstName = &s
	}
	if c.LastName != nil {
		s := utils.SanitizeString(*c.LastName)
		c.LastName = &s
	}
	if c.AvatarURL != nil {
		s := utils.SanitizeString(*c.AvatarURL)
		c.AvatarURL = &s
	}
}

func (r *RequestContact) Sanitize() {
	r.Username = utils.SanitizeString(r.Username)
}
