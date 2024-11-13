package validators

import (
	"errors"
	"regexp"
	"unicode/utf8"

	"github.com/jim-ww/nms-go/internal/dtos"
)

var (
	ErrUsernameLength = errors.New("username length must be between 3 and 30")
	ErrPasswordLength = errors.New("password length must be between 3 and 255")
	ErrEmailLength    = errors.New("email length must be between 3 and 255")
	ErrInvalidEmail   = errors.New("invalid email")
	// TODO does this belong in here?
	ErrUsernameTaken = errors.New("username already exists")
	ErrEmailTaken    = errors.New("email already exists")
)

func ValidateLoginDTO(dto *dtos.LoginDTO) (isValid bool, errs []error) {
	usernameValid, err := isValidUsername(dto.Username)
	if err != nil {
		errs = append(errs, err)
	}
	passwordValid, err := isValidPassword(dto.Password)
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return false, errs
	}
	return usernameValid && passwordValid, nil
}

func ValidateRegisterDTO(dto *dtos.RegisterDTO) (isValid bool, errs []error) {
	usernameValid, err := isValidUsername(dto.Username)
	if err != nil {
		errs = append(errs, err)
	}
	passwordValid, err := isValidPassword(dto.Password)
	if err != nil {
		errs = append(errs, err)
	}
	emailValid, err := isValidEmail(dto.Email)
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return false, errs
	}
	return usernameValid && passwordValid && emailValid, nil
}

func isValidUsername(username string) (bool, error) {
	len := utf8.RuneCountInString(username)
	if len < 3 || len > 30 {
		return false, ErrUsernameLength
	}
	return true, nil
}

func isValidEmail(email string) (bool, error) {
	len := utf8.RuneCountInString(email)
	if len < 3 || len > 255 {
		return false, ErrEmailLength
	}
	re := regexp.MustCompile(`^[\w\-\.]+@([\w-]+\.)+[\w-]{2,}$`)
	if !re.MatchString(email) {
		return false, ErrInvalidEmail
	}
	return true, nil
}

func isValidPassword(password string) (bool, error) {
	len := utf8.RuneCountInString(password)
	if len < 3 || len > 255 {
		return false, ErrPasswordLength
	}
	return true, nil
}
