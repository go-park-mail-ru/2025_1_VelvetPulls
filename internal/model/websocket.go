package model

type MessageEvent struct {
	Action  string  `json:"action"`
	Message Message `json:"payload"`
}

type ChatEvent struct {
	Action string   `json:"action"`
	Chat   ChatInfo `json:"payload"`
}
