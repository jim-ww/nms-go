package auth

import (
	"regexp"
	"unicode/utf8"

	loginDTO "github.com/jim-ww/nms-go/internal/features/auth/dtos/login"
	registerDTO "github.com/jim-ww/nms-go/internal/features/auth/dtos/register"
)

type FieldError string

const (
	ErrUsernameLength FieldError = "username length must be between 3 and 30 characters"
	ErrPasswordLength FieldError = "password length must be between 3 and 255 characters"
	ErrEmailLength    FieldError = "email length must be between 3 and 255 characters"
	ErrEmailInvalid   FieldError = "invalid email"
)

type Field string

const (
	UsernameField Field = "username"
	EmailField    Field = "email"
	PasswordField Field = "password"
)

type ValidationErrors map[Field][]FieldError

func (vErrs ValidationErrors) HasErrors() bool {
	return len(vErrs[UsernameField]) > 0 || len(vErrs[EmailField]) > 0 || len(vErrs[PasswordField]) > 0
}

func (errors ValidationErrors) TranslateValidationErrors() map[string][]string {
	translated := make(map[string][]string, 3)

	for field, errs := range errors {
		key := string(field)
		for _, err := range errs {
			translated[key] = append(translated[key], string(err)) // Convert FieldError to string
		}
	}

	return translated
}

func ValidateLoginDTO(dto *loginDTO.LoginDTO) ValidationErrors {
	errs := make(ValidationErrors, 2)
	errs[UsernameField] = validateUsername(dto.Username)
	errs[PasswordField] = validatePassword(dto.Password)
	return errs
}

func ValidateRegisterDTO(dto *registerDTO.RegisterDTO) ValidationErrors {
	errs := make(ValidationErrors, 3)
	errs[UsernameField] = validateUsername(dto.Username)
	errs[EmailField] = validateEmail(dto.Email)
	errs[PasswordField] = validatePassword(dto.Password)
	return errs
}

func validateUsername(username string) (errs []FieldError) {
	len := utf8.RuneCountInString(username)
	if len < 3 || len > 30 {
		errs = append(errs, ErrUsernameLength)
	}
	return errs
}

func validateEmail(email string) (errs []FieldError) {
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

func validatePassword(password string) (errs []FieldError) {
	len := utf8.RuneCountInString(password)
	if len < 3 || len > 255 {
		errs = append(errs, ErrPasswordLength)
	}
	return errs
}
