package usecase

import (
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"github.com/go-park-mail-ru/2025_1_VelvetPulls/config"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/model"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/internal/repository"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type IUserUsecase interface {
	UploadAvatar(ctx context.Context, uploadedData multipart.File, header *multipart.FileHeader) error
	GetUserAvatar(ctx context.Context, id string) ([]byte, error)
	avatarExists(ctx context.Context, userID string) (bool, error)
	GetUserProfile(ctx context.Context, id string) (*model.UserProfile, error)
	UpdateUserProfile(ctx context.Context, profile *model.UserProfile) error
}

type UserUsecase struct {
	userRepo repository.IUserRepo
}

func NewUserUsecase(userRepo repository.IUserRepo) IUserUsecase {
	return &UserUsecase{
		userRepo: userRepo,
	}
}

func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}

	return false
}

func (uc *UserUsecase) avatarExists(ctx context.Context, userID string) (bool, error) {
	path, err := uc.userRepo.GetAvatarPathByUserID(ctx, userID)
	if err != nil {
		return false, err
	} else if path == "" {
		return false, nil
	}

	return true, nil
}

func (uc *UserUsecase) UploadAvatar(ctx context.Context, uploadedData multipart.File, header *multipart.FileHeader) error {
	userID := utils.GetUserIDFromCtx(ctx)

	if header.Size > 2*1024*1024 {
		return ErrBadAvatarSize
	}

	mimeType := header.Header.Get("Content-Type")
	ext, err := mime.ExtensionsByType(mimeType)
	if ext == nil {
		return ErrBadAvatarType
	} else if err != nil {
		return err
	}

	if contains(ext, ".jpeg") {
		header.Filename = userID + ".jpeg"
	} else if contains(ext, ".png") {
		header.Filename = userID + ".png"
	} else if contains(ext, ".gif") {
		header.Filename = userID + ".gif"
	} else {
		return ErrBadAvatarType
	}

	avaExists, err := uc.avatarExists(ctx, userID)
	if err != nil {
		return err
	}

	if avaExists {
		path, err := uc.userRepo.GetAvatarPathByUserID(ctx, userID)
		if err != nil {
			return err
		}

		delErr := os.Remove(path)
		if delErr != nil {
			return err
		}
	}

	now := time.Now()
	year := strconv.Itoa(now.Year())
	month := strconv.Itoa(int(now.Month()))
	day := strconv.Itoa(now.Day())

	dirToSave := config.UPLOADS_DIR + year + "/" + month + "/" + day
	err = os.MkdirAll(dirToSave, 0777)
	if err != nil {
		return err
	}

	filePath := dirToSave + "/" + header.Filename
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	_, copyErr := io.Copy(f, uploadedData)
	if copyErr != nil {
		return copyErr
	}

	syncErr := f.Sync()
	if syncErr != nil {
		return syncErr
	}
	closeErr := f.Close()
	if closeErr != nil {
		return closeErr
	}

	err = uc.userRepo.UpdateAvatarPathByUserID(ctx, userID, filePath)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UserUsecase) GetUserAvatar(ctx context.Context, id string) ([]byte, error) {
	path, err := uc.userRepo.GetAvatarPathByUserID(ctx, id)

	if path == "" {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, ErrReadAvatar
	}

	return fileBytes, nil
}

func (uc *UserUsecase) GetUserProfile(ctx context.Context, id string) (*model.UserProfile, error) {
	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	profile := &model.UserProfile{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Phone:     user.Phone,
		Email:     user.Email,
	}

	return profile, nil
}

func (uc *UserUsecase) UpdateUserProfile(ctx context.Context, profile *model.UserProfile) error {
	userIDstr := utils.GetUserIDFromCtx(ctx)
	userID, err := uuid.Parse(userIDstr)
	if err != nil {
		return fmt.Errorf("failed to parse user ID: %v", err)
	}

	updatedUser := &model.User{
		ID:        userID,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		Username:  profile.Username,
		Phone:     profile.Phone,
		Email:     profile.Email,
		UpdatedAt: time.Now(),
	}
	updStatus := uc.userRepo.UpdateUser(ctx, updatedUser)
	if updStatus != nil {
		return updStatus
	}

	return nil
}
