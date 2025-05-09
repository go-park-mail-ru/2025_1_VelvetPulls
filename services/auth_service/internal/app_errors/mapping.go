package apperrors

import (
	"errors"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ValidationErrorMsg   = "validation error"
	InternalErrorMsg     = "internal server error"
	UsernameTakenMsg     = "username is already taken"
	PhoneTakenMsg        = "phone number is already taken"
	InvalidUsernameMsg   = "invalid username"
	InvalidPasswordMsg   = "invalid password"
	HashPasswordErrorMsg = "failed to hash password"
)

var ErrorToGrpcStatus = map[error]struct {
	code    codes.Code
	message string
}{
	usecase.ErrUsernameIsTaken: {
		code:    codes.AlreadyExists,
		message: UsernameTakenMsg,
	},
	usecase.ErrPhoneIsTaken: {
		code:    codes.AlreadyExists,
		message: PhoneTakenMsg,
	},
	usecase.ErrHashPassword: {
		code:    codes.Internal,
		message: HashPasswordErrorMsg,
	},
	usecase.ErrInvalidUsername: {
		code:    codes.InvalidArgument,
		message: InvalidUsernameMsg,
	},
	usecase.ErrInvalidPassword: {
		code:    codes.Unauthenticated,
		message: InvalidPasswordMsg,
	},
}

// ConvertError преобразует доменные ошибки в gRPC ошибки
func ConvertError(err error) error {
	if err == nil {
		return nil
	}

	// Если это ошибка валидации
	if errors.Is(err, model.ErrValidation) {
		return status.Error(codes.InvalidArgument, ValidationErrorMsg)
	}

	// Если ошибка замаплена
	for srcErr, grpcStatus := range ErrorToGrpcStatus {
		if errors.Is(err, srcErr) {
			return status.Error(grpcStatus.code, grpcStatus.message)
		}
	}

	// Любая другая ошибка
	return status.Error(codes.Internal, InternalErrorMsg)
}
