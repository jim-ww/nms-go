package health

import (
	"github.com/jim-ww/nms-go/internal/lib/api/response"
	"github.com/labstack/echo/v4"
)

func Health(c echo.Context) error {
	return c.JSON(200, response.OK())
}
