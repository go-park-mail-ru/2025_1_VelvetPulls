package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int64          `json:"id"`
	FirstName sql.NullString `json:"first_name"`
	LastName  sql.NullString `json:"last_name"`
	Username  string         `json:"username"`
	Phone     string         `json:"phone"`
	Email     sql.NullString `json:"email"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
