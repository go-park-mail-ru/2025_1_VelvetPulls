package models

import (
	"reflect"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	user := User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Phone:     "+1234567890",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expectedUser := User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
		Phone:     "+1234567890",
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if !reflect.DeepEqual(user, expectedUser) {
		t.Errorf("User struct is incorrect. Got: %+v, Expected: %+v", user, expectedUser)
	}
}

func TestChat(t *testing.T) {
	chat := Chat{
		ID:          1,
		Type:        ChatTypeGroup,
		Title:       "Developers",
		Description: "Chat for developers",
		Members:     []int64{1, 2, 3},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	expectedChat := Chat{
		ID:          1,
		Type:        ChatTypeGroup,
		Title:       "Developers",
		Description: "Chat for developers",
		Members:     []int64{1, 2, 3},
		CreatedAt:   chat.CreatedAt,
		UpdatedAt:   chat.UpdatedAt,
	}

	if !reflect.DeepEqual(chat, expectedChat) {
		t.Errorf("Chat struct is incorrect. Got: %+v, Expected: %+v", chat, expectedChat)
	}
}

func TestMessage(t *testing.T) {
	message := Message{
		ID:        1,
		ChatID:    1,
		UserID:    1,
		Text:      "Hello, world!",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expectedMessage := Message{
		ID:        1,
		ChatID:    1,
		UserID:    1,
		Text:      "Hello, world!",
		CreatedAt: message.CreatedAt,
		UpdatedAt: message.UpdatedAt,
	}

	if !reflect.DeepEqual(message, expectedMessage) {
		t.Errorf("Message struct is incorrect. Got: %+v, Expected: %+v", message, expectedMessage)
	}
}
