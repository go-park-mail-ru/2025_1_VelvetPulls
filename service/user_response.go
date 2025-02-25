package service

type UserResponse struct {
	StatusCode int `json:"statusCode"`
	Body       any `json:"body"`
}
