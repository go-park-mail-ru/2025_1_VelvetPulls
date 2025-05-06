package model

import (
	"errors"
	"mime/multipart"

	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2025_1_VelvetPulls/pkg/utils"
	"github.com/google/uuid"
)

type GetUserProfile struct {
	AvatarPath *string `json:"avatar_path,omitempty"`
	FirstName  *string `json:"first_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	Username   string  `json:"username"`
	Phone      string  `json:"phone"`
	Email      *string `json:"email,omitempty"`
}

type UpdateUserProfile struct {
	ID        uuid.UUID       `json:"id"`
	Avatar    *multipart.File `json:"-" valid:"-"`
	FirstName *string         `json:"first_name,omitempty" valid:"optional,stringlength(1|50)"`
	LastName  *string         `json:"last_name,omitempty" valid:"optional,stringlength(1|50)"`
	Username  *string         `json:"username" valid:"optional,alphanum,length(3|20)"`
	Phone     *string         `json:"phone" valid:"optional"`
	Email     *string         `json:"email,omitempty" valid:"optional,email"`
	Password  string          `json:"password,omitempty" valid:"optional,stringlength(8|100)"`
}

func (up *UpdateUserProfile) Validate() error {
	if up.FirstName == nil && up.LastName == nil && up.Username == nil &&
		up.Phone == nil && up.Email == nil && up.Avatar == nil && up.Password == "" {
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

func (g *GetUserProfile) Sanitize() {
	g.Username = utils.SanitizeString(g.Username)
	g.Phone = utils.SanitizeString(g.Phone)

	if g.FirstName != nil {
		s := utils.SanitizeString(*g.FirstName)
		g.FirstName = &s
	}
	if g.LastName != nil {
		s := utils.SanitizeString(*g.LastName)
		g.LastName = &s
	}
	if g.Email != nil {
		s := utils.SanitizeString(*g.Email)
		g.Email = &s
	}
}

func (u *UpdateUserProfile) Sanitize() {
	if u.FirstName != nil {
		s := utils.SanitizeString(*u.FirstName)
		u.FirstName = &s
	}
	if u.LastName != nil {
		s := utils.SanitizeString(*u.LastName)
		u.LastName = &s
	}
	if u.Username != nil {
		s := utils.SanitizeString(*u.Username)
		u.Username = &s
	}
	if u.Phone != nil {
		s := utils.SanitizeString(*u.Phone)
		u.Phone = &s
	}
	if u.Email != nil {
		s := utils.SanitizeString(*u.Email)
		u.Email = &s
	}
}
