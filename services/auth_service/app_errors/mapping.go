package apperrors

import (
	"errors"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/usecase"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/services/auth_service/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Константные сообщения об ошибках
const (
	ValidationErrorMsg   = "validation error"
	InternalErrorMsg     = "internal server error"
	UsernameTakenMsg     = "username is already taken"
	PhoneTakenMsg        = "phone number is already taken"
	InvalidUsernameMsg   = "invalid username"
	InvalidPasswordMsg   = "invalid password"
	HashPasswordErrorMsg = "failed to hash password"
)

// Маппинг ошибок с константными сообщениями
var ErrorToGrpcStatus = map[error]struct {
	code    codes.Code
	message string
}{
	usecase.ErrUsernameIsTaken: {codes.AlreadyExists, UsernameTakenMsg},
	usecase.ErrPhoneIsTaken:    {codes.AlreadyExists, PhoneTakenMsg},
	usecase.ErrHashPassword:    {codes.Internal, HashPasswordErrorMsg},
	usecase.ErrInvalidUsername: {codes.InvalidArgument, InvalidUsernameMsg},
	usecase.ErrInvalidPassword: {codes.Unauthenticated, InvalidPasswordMsg},
}

func ConvertError(err error) error {
	if err == nil {
		return nil
	}

	// Обработка ошибки валидации
	if errors.Is(err, model.ErrValidation) {
		return status.Error(codes.InvalidArgument, ValidationErrorMsg)
	}

	// Проверка замапленных ошибок
	for srcErr, grpcStatus := range ErrorToGrpcStatus {
		if errors.Is(err, srcErr) {
			return status.Error(grpcStatus.code, grpcStatus.message)
		}
	}

	// Ошибка по умолчанию
	return status.Error(codes.Internal, InternalErrorMsg)
}
