package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	AvatarPath *string   `json:"avatar_path"`
	FirstName  *string   `json:"first_name"`
	LastName   *string   `json:"last_name"`
	Username   string    `json:"username"`
	Phone      string    `json:"phone"`
	Email      *string   `json:"email"`
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
