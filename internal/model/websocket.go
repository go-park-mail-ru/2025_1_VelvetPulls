package model

import "github.com/google/uuid"

type MessageEvent struct {
	Action  string  `json:"action"`
	Message Message `json:"payload"`
}

type ChatEvent struct {
	Action string   `json:"action"`
	Chat   ChatInfo `json:"payload"`
}

type AnyEvent struct {
	TypeOfEvent string
	Event       interface{}
}

type Event struct {
	Action string      `json:"action"`
	ChatId uuid.UUID   `json:"chatId"`
	Users  []uuid.UUID `json:"users"`
}

type ChatInfoWS struct {
	Events chan Event
	Users  map[uuid.UUID]struct{}
}
