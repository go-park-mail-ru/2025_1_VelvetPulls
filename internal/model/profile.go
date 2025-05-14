package model

import (
	"errors"
	"mime/multipart"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type GetUserProfile struct {
	ID         uuid.UUID  `json:"id"`
	AvatarPath *string    `json:"avatar_path,omitempty"`
	Name       string     `json:"name"`
	Username   string     `json:"username"`
	BirthDate  *time.Time `json:"birth_date,omitempty"`
}

type UpdateUserProfile struct {
	ID        uuid.UUID       `json:"id"`
	Avatar    *multipart.File `json:"-" valid:"-"`
	Name      *string         `json:"name,omitempty" valid:"optional,stringlength(1|100)"`
	Username  *string         `json:"username,omitempty" valid:"optional,alphanum,stringlength(3|20)"`
	BirthDate *time.Time      `json:"birth_date,omitempty"`
	Password  string          `json:"password,omitempty" valid:"optional,stringlength(8|100)"`
}

// Валидация: хотя бы одно поле должно быть заполнено
func (up *UpdateUserProfile) Validate() error {
	if up.Name == nil && up.Username == nil &&
		up.Avatar == nil && up.BirthDate == nil && up.Password == "" {
		return errors.Join(ErrValidation, errors.New("at least one field must be provided for update"))
	}

	if _, err := govalidator.ValidateStruct(up); err != nil {
		if errs, ok := err.(govalidator.Errors); ok {
			return errors.Join(ErrValidation, errors.New(errs.Error()))
		}
		return errors.Join(ErrValidation, err)
	}

	return nil
}

// Санитизация
func (g *GetUserProfile) Sanitize() {
	g.Username = utils.SanitizeString(g.Username)
	g.Name = utils.SanitizeString(g.Name)
}

func (u *UpdateUserProfile) Sanitize() {
	if u.Name != nil {
		s := utils.SanitizeString(*u.Name)
		u.Name = &s
	}
	if u.Username != nil {
		s := utils.SanitizeString(*u.Username)
		u.Username = &s
	}
}
