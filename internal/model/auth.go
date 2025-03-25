package model

import (
	"errors"
	"regexp"

	"github.com/asaskevich/govalidator"
)

var passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_\-+=]{8,32}$`)

type LoginCredentials struct {
	Username string `json:"username" valid:"required,alphanum,length(3|20)"`
	Password string `json:"password" valid:"required,length(6|100)"`
}

type RegisterCredentials struct {
	Username        string `json:"username" valid:"required,alphanum,length(3|20)"`
	Password        string `json:"password" valid:"required,length(6|100)"`
	ConfirmPassword string `json:"confirm_password" valid:"required,length(6|100)"`
	Phone           string `json:"phone" valid:"required,numeric,length(10|15)"`
}

func (lc *LoginCredentials) Validate() error {
	_, err := govalidator.ValidateStruct(lc)
	if err != nil {
		return err
	}

	return nil
}

func (rc *RegisterCredentials) Validate() error {
	_, err := govalidator.ValidateStruct(rc)
	if err != nil {
		return err
	}

	if !passwordRegex.MatchString(rc.Password) {
		return errors.New("password must contain at least one special character, one letter, and one digit")
	}

	if rc.Password != rc.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	return nil
}
