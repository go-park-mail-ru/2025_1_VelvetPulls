package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID  `json:"id"`
	AvatarPath *string    `json:"avatarPath,omitempty"`
	Name       string     `json:"name"`
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	BirthDate  *time.Time `json:"birth_date,omitempty"`
}
