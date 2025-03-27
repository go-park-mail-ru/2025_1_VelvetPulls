package model

import (
	"errors"
	"regexp"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

var phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)

type User struct {
	ID         uuid.UUID `json:"id"`
	AvatarPath *string   `json:"avatar_path"`
	FirstName  *string   `json:"first_name"`
	LastName   *string   `json:"last_name"`
	Username   string    `json:"username"`
	Phone      string    `json:"phone"`
	Email      *string   `json:"email"`
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserProfile struct {
	ID        uuid.UUID `json:"id"`
	FirstName *string   `json:"first_name" valid:"optional,stringlength(1|50)"`
	LastName  *string   `json:"last_name" valid:"optional,stringlength(1|50)"`
	Username  string    `json:"username" valid:"optional,alphanum,length(3|20)"`
	Phone     string    `json:"phone" valid:"optional"`
	Email     *string   `json:"email" valid:"optional,email"`
}

func (up *UserProfile) Validate() error {
	_, err := govalidator.ValidateStruct(up)
	if err != nil {
		return err
	}

	if up.FirstName == nil && up.LastName == nil && up.Username == "" && up.Phone == "" && up.Email == nil {
		return errors.New("at least one field must be provided for update")
	}

	if up.Phone != "" && !phoneRegex.MatchString(up.Phone) {
		return errors.New("invalid phone number format")
	}

	return nil
}
