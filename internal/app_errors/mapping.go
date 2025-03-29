package apperrors

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	servererrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/server_errors"
)

var errToCode = map[error]int{
	// HTTP errors
	servererrors.ErrInvalidRequestData: http.StatusBadRequest, // 400
	servererrors.ErrValidation:         http.StatusBadRequest,
	// Usecase errors
	usecase.ErrUsernameIsTaken: http.StatusConflict,            // 409 - Конфликт, имя пользователя занято
	usecase.ErrPhoneIsTaken:    http.StatusConflict,            // 409 - Конфликт, телефон занят
	usecase.ErrHashPassword:    http.StatusInternalServerError, // 500 - Ошибка сервера при хешировании
	usecase.ErrInvalidUsername: http.StatusBadRequest,          // 400 - Некорректное имя пользователя
	usecase.ErrInvalidPassword: http.StatusBadRequest,          // 400 - Некорректный пароль

	// Repository errors
	repository.ErrSessionNotFound:     http.StatusNotFound,            // 404
	repository.ErrUserNotFound:        http.StatusNotFound,            // 404
	repository.ErrRecordAlreadyExists: http.StatusConflict,            // 409
	repository.ErrUpdateFailed:        http.StatusInternalServerError, // 500
	repository.ErrInvalidUUID:         http.StatusBadRequest,          // 400
	repository.ErrEmptyField:          http.StatusBadRequest,          // 400
	repository.ErrDatabaseOperation:   http.StatusInternalServerError, // 500
}

func GetErrAndCodeToSend(err error) (int, error) {
	var source error
	for err != nil {
		source = err
		err = errors.Unwrap(err)
	}

	code, ok := errToCode[source]
	if !ok {
		return http.StatusInternalServerError, servererrors.ErrInternalServer
	}

	return code, source
}
