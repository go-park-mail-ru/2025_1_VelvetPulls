package model

import (
	"errors"
	"mime/multipart"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type Chat struct {
	ID         uuid.UUID `json:"id" valid:"uuid"`
	AvatarPath *string   `json:"avatar_path"`
	Type       string    `json:"type" valid:"in(dialog|group|channel),required"`
	Title      string    `json:"title" valid:"required~Title is required,length(1|100)"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
}

type CreateChatRequest struct {
	Type       string `json:"type" valid:"in(dialog|group|channel),required"`
	Title      string `json:"title" valid:"required~Title is required,length(1|100)"`
	DialogUser string `json:"dialog_user,omitempty" valid:"-"`
}

type CreateChat struct {
	Avatar     *multipart.File `json:"-" valid:"-"`
	Type       string          `json:"type" valid:"in(dialog|group|channel),required"`
	Title      string          `json:"title" valid:"required~Title is required,length(1|100)"`
	DialogUser string          `json:"-" valid:"-"`
}

type UpdateChat struct {
	ID     uuid.UUID       `json:"id" valid:"required,uuid"`
	Avatar *multipart.File `json:"-" valid:"-"`
	Title  *string         `json:"title" valid:"length(1|100)"`
}

type ChatInfo struct {
	ID         uuid.UUID    `json:"id" valid:"uuid"`
	AvatarPath *string      `json:"avatar_path"`
	Type       string       `json:"type" valid:"in(dialog|group|channel)"`
	Title      string       `json:"title" valid:"length(1|100)"`
	CountUsers int          `json:"count_users" valid:"range(0|5000)"`
	Users      []UserInChat `json:"users"`
}

type UserInChat struct {
	ID         uuid.UUID `json:"id" valid:"uuid"`
	Username   string    `json:"username" valid:"required,length(3|50)"`
	Name       *string   `json:"name" valid:"length(0|100)"`
	AvatarPath *string   `json:"avatar_path"`
	Role       *string   `json:"role" valid:"length(0|20)"`
}

type AddedUsersIntoChat struct {
	AddedUsers    []string `json:"added_users" valid:"required"`
	NotAddedUsers []string `json:"not_added_users" valid:"required"`
}

type DeletedUsersFromChat struct {
	DeletedUsers []string `json:"deleted_users" valid:"required"`
}

func (c *Chat) Validate() error {
	if _, err := govalidator.ValidateStruct(c); err != nil {
		return errors.Join(ErrValidation, errors.New("invalid chat data: "+err.Error()))
	}
	return nil
}

func (c *CreateChat) Validate() error {
	if _, err := govalidator.ValidateStruct(c); err != nil {
		return errors.Join(ErrValidation, errors.New("invalid create chat data: "+err.Error()))
	}
	return nil
}

func (u *UpdateChat) Validate() error {
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return errors.Join(ErrValidation, errors.New("invalid update chat data: "+err.Error()))
	}
	return nil
}

func (c *ChatInfo) Validate() error {
	if _, err := govalidator.ValidateStruct(c); err != nil {
		return errors.Join(ErrValidation, errors.New("invalid chat info data: "+err.Error()))
	}
	return nil
}

func (u *UserInChat) Validate() error {
	if _, err := govalidator.ValidateStruct(u); err != nil {
		return errors.Join(ErrValidation, errors.New("invalid user in chat data: "+err.Error()))
	}
	return nil
}

func (a *AddedUsersIntoChat) Validate() error {
	if _, err := govalidator.ValidateStruct(a); err != nil {
		return errors.Join(ErrValidation, errors.New("invalid added users data: "+err.Error()))
	}
	return nil
}

func (d *DeletedUsersFromChat) Validate() error {
	if _, err := govalidator.ValidateStruct(d); err != nil {
		return errors.Join(ErrValidation, errors.New("invalid deleted users data: "+err.Error()))
	}
	return nil
}
