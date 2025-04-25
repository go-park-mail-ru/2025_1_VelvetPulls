package apperrors

import (
	"errors"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	servererrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/server_errors"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errToCode = map[error]int{
	// HTTP errors
	servererrors.ErrInvalidRequestData: http.StatusBadRequest, // 400

	// Usecase errors
	usecase.ErrUsernameIsTaken:         http.StatusConflict,            // 409
	usecase.ErrPhoneIsTaken:            http.StatusConflict,            // 409
	usecase.ErrHashPassword:            http.StatusInternalServerError, // 500
	usecase.ErrInvalidUsername:         http.StatusBadRequest,          // 400
	usecase.ErrInvalidPassword:         http.StatusBadRequest,          // 400
	usecase.ErrPermissionDenied:        http.StatusForbidden,
	usecase.ErrDialogUpdateForbidden:   http.StatusBadRequest,
	usecase.ErrOnlyOwnerCanModify:      http.StatusForbidden,
	usecase.ErrDialogAddUsers:          http.StatusBadRequest,
	usecase.ErrDialogDeleteUsers:       http.StatusBadRequest,
	usecase.ErrChatCreationFailed:      http.StatusInternalServerError,
	usecase.ErrAddOwnerToDialog:        http.StatusInternalServerError,
	usecase.ErrAddParticipantToDialog:  http.StatusInternalServerError,
	usecase.ErrAddOwnerToGroup:         http.StatusInternalServerError,
	usecase.ErrOnlyOwnerCanDelete:      http.StatusForbidden,
	usecase.ErrOnlyOwnerCanAddUsers:    http.StatusForbidden,
	usecase.ErrOnlyOwnerCanDeleteUsers: http.StatusForbidden,

	// Repository errors
	repository.ErrSessionNotFound:     http.StatusNotFound,            // 404
	repository.ErrSelfContact:         http.StatusBadRequest,          // 400
	repository.ErrUserNotFound:        http.StatusNotFound,            // 404
	repository.ErrChatNotFound:        http.StatusNotFound,            // 404
	repository.ErrRecordAlreadyExists: http.StatusConflict,            // 409
	repository.ErrUpdateFailed:        http.StatusInternalServerError, // 500
	repository.ErrInvalidUUID:         http.StatusBadRequest,          // 400
	repository.ErrEmptyField:          http.StatusBadRequest,          // 400
	repository.ErrDatabaseOperation:   http.StatusInternalServerError, // 500
	repository.ErrDatabaseScan:        http.StatusInternalServerError, // 500

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

func GrpcCodeToHttp(grpcCode codes.Code) int {
	switch grpcCode {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.Internal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// UnpackGrpcError извлекает сообщение и код из gRPC-ошибки
func UnpackGrpcError(err error) (int, string) {
	if err == nil {
		return http.StatusOK, ""
	}

	// Пытаемся распарсить gRPC-статус
	st, ok := status.FromError(err)
	if ok {
		return GrpcCodeToHttp(st.Code()), st.Message()
	}

	code, err := GetErrAndCodeToSend(err)
	return code, err.Error()
}
