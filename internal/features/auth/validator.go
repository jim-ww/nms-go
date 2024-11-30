package auth

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/jim-ww/nms-go/internal/features/auth/dtos"
)

type Field string
type FieldError string

const (
	UsernameField Field = "username"
	EmailField    Field = "email"
	PasswordField Field = "password"
)

const (
	ErrUsernameTaken        FieldError = "username already exists"
	ErrEmailTaken           FieldError = "email already exists"
	ErrUsernameDoesNotExist FieldError = "username does not exist"
	ErrInvalidPassword      FieldError = "invalid password"
)

type ValidationErrors *map[Field][]FieldError

type AuthValidator struct {
	logger   *slog.Logger
	validatr *validator.Validate
}

func (v *AuthValidator) ValidateLoginDTO(dto dtos.LoginDTO) (errors ValidationErrors) {
	if err := v.validatr.Struct(dto); err != nil {
		v.logger.Debug("field validation completed with errors")
	}
	return nil
}
