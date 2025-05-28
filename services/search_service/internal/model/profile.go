package model

import (
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type UserProfile struct {
	ID         uuid.UUID  `json:"id"`
	AvatarPath *string    `json:"avatar_path,omitempty"`
	Name       *string    `json:"name,omitempty"`
	Username   string     `json:"username"`
	BirthDate  *time.Time `json:"birth_date,omitempty"`
}

type RequestUserProfile struct {
	Username string `json:"username"`
}

func (g *UserProfile) Sanitize() {
	g.Username = utils.SanitizeString(g.Username)

	if g.Name != nil {
		s := utils.SanitizeString(*g.Name)
		g.Name = &s
	}
}

func (r *RequestUserProfile) Sanitize() {
	r.Username = utils.SanitizeString(r.Username)
}
