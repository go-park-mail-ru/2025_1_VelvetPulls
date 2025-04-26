package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/csat_service/model"
)

type ICsatUsecase interface {
	GetQuestion(ctx context.Context, values model.Question)
	SendAnswer(ctx context.Context, values model.Question)
	getStatistic(ctx context.Context, values model.UserInfo)
}
