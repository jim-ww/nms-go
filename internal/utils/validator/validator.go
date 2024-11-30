package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func FormatErrors(errs validator.ValidationErrors) []string {
	errorMsgs := []string{}
	for _, err := range errs {
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
	}
	return errorMsgs
}
