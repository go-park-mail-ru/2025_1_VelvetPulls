package model

type CreateCsatAnswerRequest struct {
	QuestionID string `json:"question_id" validate:"required,uuid"`
	UserID     string `json:"user_id" validate:"required,uuid"`
	Rating     int    `json:"rating" validate:"required,min=1,max=5"`
}
