package auth

import (
	"regexp"
	"unicode/utf8"
)

const (
	ErrUsernameLength = "username length must be between 3 and 30 characters"
	ErrPasswordLength = "password length must be between 3 and 255 characters"
	ErrEmailLength    = "email length must be between 3 and 255 characters"
	ErrEmailInvalid   = "invalid email"
)

const (
	UsernameField = "username"
	EmailField    = "email"
	PasswordField = "password"
)

func ValidateLoginDTO(dto *LoginDTO) map[string][]string {
	errs := make(map[string][]string, 2)
	errs[UsernameField] = validateUsername(dto.Username)
	errs[PasswordField] = validatePassword(dto.Password)
	return errs
}

func ValidateRegisterDTO(dto *RegisterDTO) map[string][]string {
	errs := make(map[string][]string, 3)
	errs[UsernameField] = validateUsername(dto.Username)
	errs[EmailField] = validateEmail(dto.Email)
	errs[PasswordField] = validatePassword(dto.Password)
	return errs
}

func validateUsername(username string) (errs []string) {
	len := utf8.RuneCountInString(username)
	if len < 3 || len > 30 {
		errs = append(errs, ErrUsernameLength)
	}
	return errs
}

func validateEmail(email string) (errs []string) {
	re := regexp.MustCompile(`^[\w\-\.]+@([\w-]+\.)+[\w-]{2,}$`)
	if !re.MatchString(email) {
		errs = append(errs, ErrEmailInvalid)
	}
	len := utf8.RuneCountInString(email)
	if len < 3 || len > 255 {
		errs = append(errs, ErrEmailLength)
	}
	return errs
}

func validatePassword(password string) (errs []string) {
	len := utf8.RuneCountInString(password)
	if len < 3 || len > 255 {
		errs = append(errs, ErrPasswordLength)
	}
	return errs
}
