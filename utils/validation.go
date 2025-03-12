package utils

import (
	"regexp"
)

// ValidatePhone проверяет правильность формата телефона для нескольких стран
func ValidateRegistrationPhone(phone string) bool {
	re := regexp.MustCompile(`^(\+7\d{10}|\+1\d{10}|\+44\d{10}|\+49\d{10})$`)
	return re.MatchString(phone)
}

// ValidateRegistrationPassword проверяет, что пароль от 8 до 32 символов,
// без пробелов и запрещённых символов (например, смайликов)
func ValidateRegistrationPassword(password string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_\-+=]{8,32}$`)
	return re.MatchString(password)
}

// ValidateRegistrationUsername проверяет, что имя пользователя:
// - Длина от 3 до 20 символов
// - Разрешены только латинские, цифры и _
func ValidateRegistrationUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	return re.MatchString(username)
}
