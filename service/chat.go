package service

import (
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/repository"
)

func FetchChatsBySession(token string) ([]*model.Chat, error) {
	session, err := repository.GetSessionBySessId(token)
	if err != nil {
		return nil, err
	}

	chats, err := repository.GetChatsByUsername(session.Username)
	if err != nil {
		return nil, err
	}

	return chats, nil
}
