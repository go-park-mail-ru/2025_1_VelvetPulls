package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
)

type ISessionUsecase interface {
	CheckLogin(ctx context.Context, token string) (string, error)
}

type SessionUsecase struct {
	sessionRepo repository.ISessionRepo
}

func NewSessionUsecase(sessionRepo repository.ISessionRepo) ISessionUsecase {
	return &SessionUsecase{
		sessionRepo: sessionRepo,
	}
}

func (uc *SessionUsecase) CheckLogin(ctx context.Context, token string) (string, error) {
	logger := utils.GetLoggerFromCtx(ctx)
	logger.Info("Checking login status")

	userID, err := uc.sessionRepo.GetUserIDByToken(ctx, token)
	if err != nil {
		return "", err
	}

	return userID, nil
}
