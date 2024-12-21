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

func formatFieldError(err validator.FieldError) []string {
	errorMsgs := []string{}
	switch err.ActualTag() {
	case "required":
		errorMsgs = append(errorMsgs, fmt.Sprintf("%s is required", err.Field()))
	case "email":
		errorMsgs = append(errorMsgs, fmt.Sprintf("'%s' is not a valid email", err.Value()))
	case "min":
		errorMsgs = append(errorMsgs, fmt.Sprintf("%s length must be greater or equal to %s", err.Field(), err.Param()))
	case "max":
		errorMsgs = append(errorMsgs, fmt.Sprintf("%s length must be less or equal to %s", err.Field(), err.Param()))
	default:
		errorMsgs = append(errorMsgs, fmt.Sprintf("%s is not valid", err.Field()))
	}
	return errorMsgs
}

func ConvertValidationErrorsToMap(validationErrors validator.ValidationErrors) map[string][]string {
	vErrs := make(map[string][]string, 3)
	for _, vErr := range validationErrors {
		vErrs[vErr.Field()] = append(vErrs[vErr.Field()], formatFieldError(vErr)...)
	}
	return vErrs
}
