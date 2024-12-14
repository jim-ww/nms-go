package errorhandler

import (
	"net/http"

	"github.com/jim-ww/nms-go/internal/lib/api/response"
	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error("Error:", err)

	switch err.(type) {
	default:
		c.JSON(http.StatusInternalServerError, response.Error("Internal Server Error"))
	}
}
