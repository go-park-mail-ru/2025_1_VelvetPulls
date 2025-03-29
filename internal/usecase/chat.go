package usecase

import (
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
)

type IChatUsecase interface {
	FetchChatsBySession(token string) ([]model.Chat, error)
}

type ChatUsecase struct {
	sessionRepo repository.ISessionRepo
	chatRepo    repository.IChatRepo
}

// NewChatUsecase создает новый экземпляр ChatUsecase.
func NewChatUsecase(sessionRepo repository.ISessionRepo, chatRepo repository.IChatRepo) IChatUsecase {
	return &ChatUsecase{
		sessionRepo: sessionRepo,
		chatRepo:    chatRepo,
	}
}

// FetchChatsBySession получает список чатов по токену сессии.
func (uc *ChatUsecase) FetchChatsBySession(token string) ([]model.Chat, error) {
	session, err := uc.sessionRepo.GetSessionBySessId(token)
	if err != nil {
		return nil, err
	}

	chats, err := uc.chatRepo.GetChatsByUsername(session)
	if err != nil {
		return nil, err
	}

	return chats, nil
}
