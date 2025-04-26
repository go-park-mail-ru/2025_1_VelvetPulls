package model

type Question struct {
	questionId string `json:"question_id"`
	askText    string `json:"ask_text"`
}

type Answer struct {
	answerId   string `json:"answer_id"`
	username   string `json:"username"`
	questionId string `json:"question_id"`
	ansText    string `json:"ans_text"`
}

type UserInfo struct {
	score    string `json:"score"`
	username string `json:"username"`
}
