package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	AvatarPath *string   `json:"avatar_path,omitempty"`
	Name       string    `json:"name"`
	BirthDate  time.Time `json:"birth_date,omitempty"`
	Username   string    `json:"username"`
	Password   string    `json:"password,omitempty"`
}
