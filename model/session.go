package model

import "time"

type Session struct {
	Username string    `json:"username"`
	Expiry   time.Time `json:"expiry"`
}

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
