package service

import (
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/repository"
)

func FetchChatsBySession(token string) (UserResponse, error) {
	session, err := repository.GetSessionBySessId(token)
	if err != nil {
		return UserResponse{
			StatusCode: 500,
			Body:       err,
		}, err
	}

	chats, err := repository.GetChatsByUsername(session.Username)
	if err != nil {
		return UserResponse{
			StatusCode: 500,
			Body:       err,
		}, err
	}

	return UserResponse{
		StatusCode: 200,
		Body:       chats,
	}, nil
}
