package model

import "time"

type Session struct {
	ID     string    `json:"username"`
	Expiry time.Time `json:"expiry"`
}

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

type AuthCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
