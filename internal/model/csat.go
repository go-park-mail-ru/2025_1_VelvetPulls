package model

type CreateCsatAnswerRequest struct {
	QuestionID string `json:"question_id" validate:"required,uuid"`
	Username   string `json:"username" validate:"required"`
	Rating     int    `json:"rating" validate:"required,min=1,max=5"`
	Feedback   string `json:"feedback"`
}
