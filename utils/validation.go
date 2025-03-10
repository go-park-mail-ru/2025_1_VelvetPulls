package utils

import (
	"regexp"
)

// validatePhone проверяет правильность формата телефона
func ValidateRegistrationPhone(phone string) bool {
	re := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	return re.MatchString(phone)
}

// validatePassword проверяет, что пароль достаточно длинный
func ValidateRegistrationPassword(password string) bool {
	return len(password) >= 8
}

// validateUsername проверяет, что имя пользователя не пустое
func ValidateRegistrationUsername(username string) bool {
	return len(username) >= 3
}
