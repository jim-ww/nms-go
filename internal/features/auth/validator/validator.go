package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/jim-ww/nms-go/internal/features/auth/dtos"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func ValidateLoginDTO(dto *dtos.LoginDTO) (errs validator.ValidationErrors) {
	if err := validate.Struct(dto); err != nil {
		validationErrors, _ := err.(validator.ValidationErrors)
		fmt.Println(validationErrors)
		return validationErrors
	}
	return errs
}

func ValidateRegisterDTO(dto *dtos.RegisterDTO) (errs validator.ValidationErrors) {
	if err := validate.Struct(dto); err != nil {
		validationErrors, _ := err.(validator.ValidationErrors)
		fmt.Println(validationErrors)
		return validationErrors
	}
	return errs
}
