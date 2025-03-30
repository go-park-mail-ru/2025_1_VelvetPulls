package apperrors

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	servererrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/server_errors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
)

var errToCode = map[error]int{
	// HTTP errors
	servererrors.ErrInvalidRequestData: http.StatusBadRequest, // 400

	// Usecase errors
	usecase.ErrUsernameIsTaken: http.StatusConflict,            // 409
	usecase.ErrPhoneIsTaken:    http.StatusConflict,            // 409
	usecase.ErrHashPassword:    http.StatusInternalServerError, // 500
	usecase.ErrInvalidUsername: http.StatusBadRequest,          // 400
	usecase.ErrInvalidPassword: http.StatusBadRequest,          // 400

	// Repository errors
	repository.ErrSessionNotFound:     http.StatusNotFound,            // 404
	repository.ErrUserNotFound:        http.StatusNotFound,            // 404
	repository.ErrRecordAlreadyExists: http.StatusConflict,            // 409
	repository.ErrUpdateFailed:        http.StatusInternalServerError, // 500
	repository.ErrInvalidUUID:         http.StatusBadRequest,          // 400
	repository.ErrEmptyField:          http.StatusBadRequest,          // 400
	repository.ErrDatabaseOperation:   http.StatusInternalServerError, // 500

	utils.ErrNotImage:      http.StatusBadRequest,          // 400
	utils.ErrSavingImage:   http.StatusInternalServerError, // 500
	utils.ErrUpdatingImage: http.StatusInternalServerError, // 500
	utils.ErrDeletingImage: http.StatusInternalServerError, // 500
}

func GetErrAndCodeToSend(err error) (int, error) {
	var source error
	for err != nil {
		if errors.Is(err, model.ErrValidation) {
			return http.StatusBadRequest, err
		}
		source = err
		err = errors.Unwrap(err)
	}

	code, ok := errToCode[source]
	if !ok {
		return http.StatusInternalServerError, servererrors.ErrInternalServer
	}
	return code, source
}
