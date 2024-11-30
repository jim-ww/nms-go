package getAuthform

import (
	"log/slog"

	"github.com/jim-ww/nms-go/internal/templates"
	"github.com/labstack/echo/v4"
)

type AuthForm struct {
	log *slog.Logger
}

func New(log *slog.Logger) *AuthForm {
	return &AuthForm{
		log: log,
	}
}

func (af AuthForm) Login(c echo.Context) error {
	tmpl := templates.Layout("Login", templates.AuthForm("", "", false, map[string][]string{}))
	return tmpl.Render(c.Request().Context(), c.Response().Writer)
}

func (af AuthForm) Register(c echo.Context) error {
	tmpl := templates.Layout("Register", templates.AuthForm("", "", true, map[string][]string{}))
	return tmpl.Render(c.Request().Context(), c.Response().Writer)
}
