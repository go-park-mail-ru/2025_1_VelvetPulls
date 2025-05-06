package model

import (
	"errors"

	"github.com/asaskevich/govalidator"
)

type Validator interface {
	Validate() error
}

type LoginCredentials struct {
	Username string `json:"username" valid:"required,length(3|20),matches(^[a-zA-Z0-9!@#$%^&*()_\\-+=]+$)"`
	Password string `json:"password" valid:"required,length(8|32),matches(^[a-zA-Z0-9!@#$%^&*()_\\-+=]+$)"`
}

type RegisterCredentials struct {
	Username        string `json:"username" valid:"required,length(3|20),matches(^[a-zA-Z0-9!@#$%^&*()_\\-+=]+$)"`
	Password        string `json:"password" valid:"required,length(8|32),matches(^[a-zA-Z0-9!@#$%^&*()_\\-+=]+$)"`
	ConfirmPassword string `json:"confirm_password" valid:"required,length(8|32)"`
	Phone           string `json:"phone" valid:"required,numeric,length(10|15)"`
}

func (lc *LoginCredentials) Validate() error {
	if _, err := govalidator.ValidateStruct(lc); err != nil {
		if errs, ok := err.(govalidator.Errors); ok {
			return errors.Join(ErrValidation, errors.New(errs.Error()))
		}
		return errors.Join(ErrValidation, err)
	}
	return nil
}

func (rc *RegisterCredentials) Validate() error {
	if _, err := govalidator.ValidateStruct(rc); err != nil {
		if errs, ok := err.(govalidator.Errors); ok {
			return errors.Join(ErrValidation, errors.New(errs.Error()))
		}
		return errors.Join(ErrValidation, err)
	}
	if rc.Password != rc.ConfirmPassword {
		return errors.Join(ErrValidation, errors.New("passwords do not match"))
	}
	return nil
}
