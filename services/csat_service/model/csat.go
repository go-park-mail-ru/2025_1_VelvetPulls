package model

import (
	"time"

	"github.com/google/uuid"
)

type RatingScale int

const (
	Rating1 RatingScale = 1
	Rating2 RatingScale = 2
	Rating3 RatingScale = 3
	Rating4 RatingScale = 4
	Rating5 RatingScale = 5
)

type Question struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	QuestionText string    `json:"question_text"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Answer struct {
	ID         uuid.UUID   `json:"id"`
	QuestionID uuid.UUID   `json:"question_id"`
	UserID     uuid.UUID   `json:"user_id"`
	Rating     RatingScale `json:"rating"`
	Feedback   *string     `json:"feedback,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
}

type UserActivity struct {
	UserID         uuid.UUID `json:"user_id"`
	LastResponseAt time.Time `json:"last_response_at"`
	ResponsesCount int       `json:"responses_count"`
}

type RatingDistribution struct {
	Rating RatingScale `json:"rating"`
	Count  int         `json:"count"`
}

type QuestionStatistics struct {
	QuestionID     uuid.UUID            `json:"question_id"`
	QuestionText   string               `json:"question_text"`
	AverageRating  float64              `json:"average_rating"`
	TotalResponses int                  `json:"total_responses"`
	Distribution   []RatingDistribution `json:"distribution"`
	Comments       []*Answer            `json:"comments"`
}

type FullStatistics struct {
	Questions []QuestionStatistics `json:"questions"`
}
