//go:generate easyjson -all user.go
package model

import (
	"time"

	"github.com/google/uuid"
)

//easyjson:json
type User struct {
	ID         uuid.UUID  `json:"id"`
	AvatarPath *string    `json:"avatar_path,omitempty"`
	Name       string     `json:"name"`
	Username   string     `json:"username"`
	Password   string     `json:"password,omitempty"`
	BirthDate  *time.Time `json:"birth_date,omitempty"`
}
