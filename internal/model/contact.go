package model

import "github.com/google/uuid"

type Contact struct {
	ID        uuid.UUID `json:"id"`
	FirstName *string   `json:"first_name"`
	LastName  *string   `json:"last_name"`
	Username  string    `json:"username"`
	AvatarURL *string   `json:"avatar_path"`
}

type RequestContact struct {
	ID uuid.UUID `json:"id"`
}
