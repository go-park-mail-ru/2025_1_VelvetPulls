package chat

import (
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
)

type ChatUsecaseInterface interface {
	FetchChatsBySession(token string) (*[]model.Chat, error)
}

type ChatUsecase struct {
	sessionRepo repository.SessionRepoInterface
	chatRepo    repository.ChatRepoInterface
}

// NewChatUsecase создает новый экземпляр ChatUsecase.
func NewChatUsecase(sessionRepo repository.SessionRepoInterface, chatRepo repository.ChatRepoInterface) *ChatUsecase {
	return &ChatUsecase{
		sessionRepo: sessionRepo,
		chatRepo:    chatRepo,
	}
}

// FetchChatsBySession получает список чатов по токену сессии.
func (uc *ChatUsecase) FetchChatsBySession(token string) (*[]model.Chat, error) {
	session, err := uc.sessionRepo.GetSessionBySessId(token)
	if err != nil {
		return nil, err
	}

	chats, err := uc.chatRepo.GetChatsByUsername(session.Username)
	if err != nil {
		return nil, err
	}

	return &chats, nil
}
