package apperrors

import (
	"errors"

	repoErrors "github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
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

	EmptyFieldMsg          = "empty required field"
	InvalidInputMsg        = "invalid input"
	RecordAlreadyExistsMsg = "record already exists"
	NotFoundMsg            = "record not found"
	DatabaseErrorMsg       = "database operation failed"
	InvalidUUIDMsg         = "invalid UUID format"
	InvalidRatingMsg       = "invalid rating value"
	NoStatisticsMsg        = "no statistics available"
	UserNotActiveMsg       = "user has no activity records"
)

// Маппинг ошибок на gRPC коды
var ErrorToGrpcStatus = map[error]struct {
	code    codes.Code
	message string
}{
	usecase.ErrUsernameIsTaken: {codes.AlreadyExists, UsernameTakenMsg},
	usecase.ErrPhoneIsTaken:    {codes.AlreadyExists, PhoneTakenMsg},
	usecase.ErrHashPassword:    {codes.Internal, HashPasswordErrorMsg},
	usecase.ErrInvalidUsername: {codes.InvalidArgument, InvalidUsernameMsg},
	usecase.ErrInvalidPassword: {codes.Unauthenticated, InvalidPasswordMsg},

	repoErrors.ErrEmptyField:          {codes.InvalidArgument, EmptyFieldMsg},
	repoErrors.ErrInvalidInput:        {codes.InvalidArgument, InvalidInputMsg},
	repoErrors.ErrRecordAlreadyExists: {codes.AlreadyExists, RecordAlreadyExistsMsg},
	repoErrors.ErrDatabaseOperation:   {codes.Internal, DatabaseErrorMsg},
	repoErrors.ErrInvalidUUID:         {codes.InvalidArgument, InvalidUUIDMsg},
}

// Конвертация ошибок в gRPC статус
func ConvertError(err error) error {
	if err == nil {
		return nil
	}

	// Специальная обработка ошибки валидации
	if errors.Is(err, model.ErrValidation) {
		return status.Error(codes.InvalidArgument, ValidationErrorMsg)
	}

	// Поиск среди известных ошибок
	for srcErr, grpcStatus := range ErrorToGrpcStatus {
		if errors.Is(err, srcErr) {
			return status.Error(grpcStatus.code, grpcStatus.message)
		}
	}

	// Ошибка по умолчанию
	return status.Error(codes.Internal, InternalErrorMsg)
}
