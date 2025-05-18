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
	// HTTP/server level
	servererrors.ErrInvalidRequestData: http.StatusBadRequest, // 400

	// Usecase level
	usecase.ErrUsernameIsTaken: http.StatusConflict,            // 409
	usecase.ErrPhoneIsTaken:    http.StatusConflict,            // 409
	usecase.ErrHashPassword:    http.StatusInternalServerError, // 500
	usecase.ErrInvalidUsername: http.StatusBadRequest,          // 400
	usecase.ErrInvalidPassword: http.StatusBadRequest,          // 400

	usecase.ErrPermissionDenied:        http.StatusForbidden,           // 403
	usecase.ErrDialogUpdateForbidden:   http.StatusBadRequest,          // 400
	usecase.ErrOnlyOwnerCanModify:      http.StatusForbidden,           // 403
	usecase.ErrDialogAddUsers:          http.StatusBadRequest,          // 400
	usecase.ErrDialogDeleteUsers:       http.StatusBadRequest,          // 400
	usecase.ErrChatCreationFailed:      http.StatusInternalServerError, // 500
	usecase.ErrAddOwnerToDialog:        http.StatusInternalServerError, // 500
	usecase.ErrAddParticipantToDialog:  http.StatusInternalServerError, // 500
	usecase.ErrAddOwnerToGroup:         http.StatusInternalServerError, // 500
	usecase.ErrOnlyOwnerCanDelete:      http.StatusForbidden,           // 403
	usecase.ErrOnlyOwnerCanAddUsers:    http.StatusForbidden,           // 403
	usecase.ErrOnlyOwnerCanDeleteUsers: http.StatusForbidden,           // 403

	usecase.ErrMessageValidationFailed: http.StatusBadRequest,          // 400
	usecase.ErrMessageCreationFailed:   http.StatusInternalServerError, // 500
	usecase.ErrMessageNotFound:         http.StatusNotFound,            // 404
	usecase.ErrMessageAccessDenied:     http.StatusForbidden,           // 403
	usecase.ErrMessageUpdateFailed:     http.StatusInternalServerError, // 500
	usecase.ErrMessageDeleteFailed:     http.StatusInternalServerError, // 500
	usecase.ErrMessagePublishFailed:    http.StatusInternalServerError, // 500

	// Repository level
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
	repository.ErrSetNotifications:    http.StatusInternalServerError, // 500

	// Utils level
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

	st, ok := status.FromError(err)
	if ok {
		return GrpcCodeToHttp(st.Code()), st.Message()
	}

	code, err := GetErrAndCodeToSend(err)
	return code, err.Error()
}
