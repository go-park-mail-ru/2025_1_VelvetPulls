package utils_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestInitLogger_NotNil(t *testing.T) {
	logger := utils.InitLogger()
	assert.NotNil(t, logger)
}

func TestInitLoggerWithFile_WritesToFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testlog-*.log")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	utils.InitLoggerWithFile(tmpFile)

	utils.Logger.Info("test message")

	// Нужно закрыть логгер, чтобы сбросить буфер
	_ = utils.Logger.Sync()
	_ = tmpFile.Sync()

	content, err := os.ReadFile(tmpFile.Name())
	assert.NoError(t, err)
	assert.Contains(t, string(content), "test message")
}

func TestLoggerWritesToStdout(t *testing.T) {
	// Перехватываем stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Переинициализируем логгер, чтобы писать в новый stdout
	logger := utils.InitLogger()
	logger.Info("stdout log test")
	_ = logger.Sync()

	// Завершаем перехват
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	assert.Contains(t, output, "stdout log test")
}
