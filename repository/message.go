package repository

import (
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/apperrors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/models"
)

var messages = []models.Message{
	{
		ID:        1,
		ChatID:    1,
		UserID:    1,
		Text:      "Привет. Как дела?",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        1,
		ChatID:    1,
		UserID:    2,
		Text:      "Привет. Да вот проект пытаюсь доделать.",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        1,
		ChatID:    1,
		UserID:    1,
		Text:      "И я... Надеюсь, успеем выполнить требования к РК.",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		ID:        1,
		ChatID:    1,
		UserID:    2,
		Text:      "Ага",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}

func GetMessageByID(messageid int64) (models.Message, error) {
	for _, message := range messages {
		if message.ID == messageid {
			return message, nil
		}
	}

	return models.Message{}, apperrors.ErrMessageNotFound
}

func GetMessagesByChat(chatid int64) []models.Message {
	result := make([]models.Message, 0)

	for _, message := range messages {
		if message.ChatID == chatid {
			result = append(result, message)
		}
	}

	return result
}
