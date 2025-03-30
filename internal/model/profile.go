package model

import (
	"errors"
	"mime/multipart"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type GetUserProfile struct {
	AvatarPath *string `json:"avatar_path"`
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	Username   string  `json:"username"`
	Phone      string  `json:"phone"`
	Email      *string `json:"email"`
}

type UpdateUserProfile struct {
	ID        uuid.UUID       `json:"id"`
	Avatar    *multipart.File `json:"-" valid:"-"`
	FirstName *string         `json:"first_name" valid:"optional,stringlength(1|50)"`
	LastName  *string         `json:"last_name" valid:"optional,stringlength(1|50)"`
	Username  *string         `json:"username" valid:"optional,alphanum,length(3|20)"`
	Phone     *string         `json:"phone" valid:"optional"`
	Email     *string         `json:"email" valid:"optional,email"`
}

func (up *UpdateUserProfile) Validate() error {
	if up.FirstName == nil && up.LastName == nil && up.Username == nil &&
		up.Phone == nil && up.Email == nil && up.Avatar == nil {
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
