package usecase

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
)

type IChatUsecase interface {
	FetchChatsByUserID(ctx context.Context) ([]model.Chat, error)
}

type ChatUsecase struct {
	sessionRepo repository.ISessionRepo
	chatRepo    repository.IChatRepo
}

func NewChatUsecase(sessionRepo repository.ISessionRepo, chatRepo repository.IChatRepo) IChatUsecase {
	return &ChatUsecase{
		sessionRepo: sessionRepo,
		chatRepo:    chatRepo,
	}
}

func (uc *ChatUsecase) FetchChatsByUserID(ctx context.Context) ([]model.Chat, error) {
	chats, err := uc.chatRepo.GetChatsByUserID(ctx)
	if err != nil {
		return nil, err
	}

	return chats, nil
}
