package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID      `json:"id"`
	AvatarPath sql.NullString `json:"avatar_path"`
	FirstName  sql.NullString `json:"first_name"`
	LastName   sql.NullString `json:"last_name"`
	Username   string         `json:"username"`
	Phone      string         `json:"phone"`
	Email      sql.NullString `json:"email"`
	Password   string         `json:"password"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
