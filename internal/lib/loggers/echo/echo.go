package echo

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func SetLevel(l echo.Logger, lvl string) {
	switch lvl {
	case "local":
		l.SetLevel(log.DEBUG)
	case "dev":
		l.SetLevel(log.DEBUG)
	case "prod":
		l.SetLevel(log.INFO)
	default:
		l.SetLevel(log.INFO)
	}
}
