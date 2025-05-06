package utils_test

import (
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestHashAndSalt_Success(t *testing.T) {
	password := "supersecure123"

	hash, err := utils.HashAndSalt(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	ok := utils.CheckPassword(hash, password)
	assert.True(t, ok)
}

func TestCheckPassword_WrongPassword(t *testing.T) {
	password := "mypassword"
	wrongPassword := "notmypassword"

	hash, err := utils.HashAndSalt(password)
	assert.NoError(t, err)

	ok := utils.CheckPassword(hash, wrongPassword)
	assert.False(t, ok)
}

func TestCheckPassword_InvalidHash(t *testing.T) {
	// строка, которая не является валидным bcrypt-хэшем
	invalidHash := "$2a$10$thisisnotavalidhash......"

	ok := utils.CheckPassword(invalidHash, "whatever")
	assert.False(t, ok)
}
