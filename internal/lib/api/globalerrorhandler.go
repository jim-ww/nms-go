package api

import (
	"net/http"

	"github.com/jim-ww/nms-go/internal/lib/api/response"
	"github.com/jim-ww/nms-go/internal/lib/errors"
	"github.com/labstack/echo/v4"
)

func GlobalErrorHandler(err error, c echo.Context) {
	c.Logger().Debug(err)
	code := http.StatusInternalServerError
	response := response.Error("Internal Server Error")

	// Check if the error is an HTTP error
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		response.Message = he.Message.(string)
	}

	switch err {
	case errors.ErrUnauthorized:
		response.Message = errors.ErrUnauthorized.Error()
		c.JSON(http.StatusForbidden, response)
	case errors.ErrInvalidJWT, errors.ErrTokenExpired, errors.ErrUnknownClaims:
		c.Redirect(http.StatusSeeOther, "/logout")
	default:
		c.JSON(code, response)
	}
}
