package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config/metrics"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/search_service/repository"
	"github.com/google/uuid"
)

type MessageUsecase struct {
	repo repository.MessageRepo
}

func NewMessageUsecase(repo repository.MessageRepo) *MessageUsecase {
	return &MessageUsecase{repo: repo}
}

func (uc *MessageUsecase) SearchMessages(
	ctx context.Context,
	chatIDStr string,
	query string,
	limit int,
	offset int,
) ([]model.Message, int, error) {
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		return nil, 0, model.ErrValidation
	}

	metrics.IncBusinessOp("search_messages")
	return uc.repo.SearchMessages(ctx, chatID, query, limit, offset)
}
