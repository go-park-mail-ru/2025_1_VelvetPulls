package model

import (
	"errors"
	"mime/multipart"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type ChatType string

const (
	ChatTypeDialog  ChatType = "dialog"
	ChatTypeGroup   ChatType = "group"
	ChatTypeChannel ChatType = "channel"
)

type UserRoleInChat string

func (ct ChatType) IsValid() bool {
	switch ct {
	case ChatTypeDialog, ChatTypeGroup:
		return true
	default:
		return false
	}
}

func (r UserRoleInChat) IsMember() bool {
	return r == RoleMember || r == RoleOwner
}

const (
	RoleOwner  UserRoleInChat = "owner"
	RoleMember UserRoleInChat = "member"
)

type Chat struct {
	ID          uuid.UUID    `json:"id" valid:"uuid"`
	AvatarPath  *string      `json:"avatar_path,omitempty"`
	Type        string       `json:"type" valid:"in(dialog|group|channel),required"`
	Title       string       `json:"title" valid:"required~Title is required,length(1|100)"`
	LastMessage *LastMessage `json:"last_message,omitempty"`
	CountUsers  int          `json:"count_users" valid:"range(0|5000)"`
}

type CreateChatRequest struct {
	Type  string `json:"type" valid:"in(dialog|group|channel),required"`
	Title string `json:"title" valid:"required~Title is required,length(1|100)"`
	//DialogUser string `json:"dialog_user,omitempty" valid:"-"`
	Users []string `json:"usersToAdd"`
}

type CreateChat struct {
	Avatar *multipart.File `json:"-" valid:"-"`
	Type   string          `json:"type" valid:"in(dialog|group|channel),required"`
	Title  string          `json:"title" valid:"required~Title is required,length(1|100)"`
	//DialogUser string          `json:"-" valid:"-"`
	Users []string `json:"usersToAdd"`
}

type UpdateChat struct {
	ID     uuid.UUID       `json:"id" valid:"required,uuid"`
	Avatar *multipart.File `json:"-" valid:"-"`
	Title  *string         `json:"title" valid:"length(1|100)"`
}

type UpdateChatResp struct {
	Avatar string `json:"updated_avatar" valid:"-"`
	Title  string `json:"title" valid:"length(1|100)"`
}

type ChatInfo struct {
	Role     string       `json:"role" example:"owner" valid:"in(owner|member)"`
	Users    []UserInChat `json:"users" valid:"-"`
	Messages []Message    `json:"messages" valid:"-"`
}

type UserInChat struct {
	ID         uuid.UUID `json:"id" valid:"uuid"`
	Username   string    `json:"username,omitempty" valid:"required,length(3|50)"`
	Name       *string   `json:"name,omitempty" valid:"length(0|100)"`
	AvatarPath *string   `json:"avatar_path,omitempty"`
	Role       *string   `json:"role" valid:"in(owner|member)"`
}

type AddedUsersIntoChat struct {
	AddedUsers    []string `json:"added_users,omitempty" valid:"required"`
	NotAddedUsers []string `json:"not_added_users,omitempty" valid:"required"`
}

type DeletedUsersFromChat struct {
	DeletedUsers []string `json:"deleted_users,omitempty" valid:"required"`
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

func (c *CreateChatRequest) Validate() error {
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

func (c *Chat) Sanitize() {
	c.Title = utils.SanitizeString(c.Title)
}

func (c *CreateChatRequest) Sanitize() {
	c.Title = utils.SanitizeString(c.Title)
	for _, user := range c.Users {
		utils.SanitizeString(user)
	}
}

func (c *CreateChat) Sanitize() {
	c.Title = utils.SanitizeString(c.Title)
	for _, user := range c.Users {
		utils.SanitizeString(user)
	}
}

func (u *UpdateChat) Sanitize() {
	if u.Title != nil {
		s := utils.SanitizeString(*u.Title)
		u.Title = &s
	}
}

func (u *UserInChat) Sanitize() {
	u.Username = utils.SanitizeString(u.Username)
	if u.Name != nil {
		s := utils.SanitizeString(*u.Name)
		u.Name = &s
	}
	if u.Role != nil {
		s := utils.SanitizeString(*u.Role)
		u.Role = &s
	}
}

func (a *AddedUsersIntoChat) Sanitize() {
	for i := range a.AddedUsers {
		a.AddedUsers[i] = utils.SanitizeString(a.AddedUsers[i])
	}
	for i := range a.NotAddedUsers {
		a.NotAddedUsers[i] = utils.SanitizeString(a.NotAddedUsers[i])
	}
}

func (d *DeletedUsersFromChat) Sanitize() {
	for i := range d.DeletedUsers {
		d.DeletedUsers[i] = utils.SanitizeString(d.DeletedUsers[i])
	}
}
