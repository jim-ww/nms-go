package auth

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/jim-ww/nms-go/internal/features/auth/dtos"
)

type Field string

const (
	UsernameField Field = "username"
	EmailField    Field = "email"
	PasswordField Field = "password"
)

const (
	UsernameTaken        = "username already exists"
	EmailTaken           = "email already exists"
	UsernameDoesNotExist = "username does not exist"
	InvalidPassword      = "invalid password"
)

type ValidationErrors map[Field][]string

func (v ValidationErrors) Len() int {
	return len(v)
}

func (v ValidationErrors) HasErrors() bool {
	return len(v[UsernameField]) > 0 || len(v[EmailField]) > 0 || len(v[PasswordField]) > 0
}

func (v ValidationErrors) TranslateToMap() map[string][]string {
	converted := make(map[string][]string, v.Len())
	for field, messages := range v {
		converted[string(field)] = messages
	}
	return converted
}

type AuthValidator struct {
	logger   *slog.Logger
	validatr *validator.Validate
}

func New(logger *slog.Logger, validatr *validator.Validate) *AuthValidator {
	return &AuthValidator{
		logger:   logger,
		validatr: validatr,
	}
}

func (v *AuthValidator) ValidateLoginDTO(dto *dtos.LoginDTO) (errors ValidationErrors) {
	if err := v.validatr.Struct(dto); err != nil {
		v.logger.Debug("field validation completed with errors")
	}
	return nil
}

func (v *AuthValidator) ValidateRegisterDTO(dto *dtos.RegisterDTO) (errors ValidationErrors) {
	if err := v.validatr.Struct(dto); err != nil {
		v.logger.Debug("field validation completed with errors")
	}
	return nil
}
