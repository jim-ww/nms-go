package handler

import (
	"fmt"
	"net/http"

	"github.com/jim-ww/nms-go/internal/features/auth/dtos"
	"github.com/jim-ww/nms-go/internal/features/auth/services/auth"
	"github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	"github.com/jim-ww/nms-go/internal/templates"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *auth.AuthService
	jwtService  *jwt.JWTService
}

func New(userService *auth.AuthService, jwtService *jwt.JWTService) *AuthHandler {
	return &AuthHandler{
		authService: userService,
		jwtService:  jwtService,
	}
}

func (AuthHandler) LoginForm(c echo.Context) error {
	tmpl := templates.Layout("Login", templates.AuthForm("", "", false, map[string][]string{}))
	return tmpl.Render(c.Request().Context(), c.Response().Writer)
}

func (AuthHandler) RegisterForm(c echo.Context) error {
	tmpl := templates.Layout("Register", templates.AuthForm("", "", true, map[string][]string{}))
	return tmpl.Render(c.Request().Context(), c.Response().Writer)
}

func (a *AuthHandler) Login(c echo.Context) error {
	req := c.Request()
	if err := req.ParseForm(); err != nil {
		return fmt.Errorf("unable to parse form %w", err)
	}
	dto := dtos.NewLoginDTO(req.FormValue("Username"), req.FormValue("Password"))
	c.Logger().Debug("got loginDTO", "dto", dto)

	token, validationErrors, err := a.authService.LoginUser(c.Request().Context(), dto)
	if err != nil {
		return fmt.Errorf("unable to login user %w", err)
	}

	if len(validationErrors) > 0 {
		c.Logger().Debug("loginDTO contains validationErrors, returning them back")
		authFormWithValidationErrors := templates.AuthForm(dto.Username, "", false, validationErrors)
		return authFormWithValidationErrors.Render(c.Request().Context(), c.Response().Writer)
	}

	c.SetCookie(a.jwtService.NewTokenCookie(token))
	c.Redirect(http.StatusSeeOther, "/")
	return nil
}

func (a *AuthHandler) Register(c echo.Context) error {
	req := c.Request()
	if err := req.ParseForm(); err != nil {
		return fmt.Errorf("unable to parse form: %w", err)
	}
	dto := dtos.NewRegisterDTO(req.FormValue("Username"), req.FormValue("Email"), req.FormValue("Password"))
	c.Logger().Debug("got registerDTO", "dto", dto)

	token, validationErrors, err := a.authService.RegisterUser(c.Request().Context(), dto)
	if err != nil {
		return fmt.Errorf("unable to register user: %w", err)
	}

	if len(validationErrors) > 0 {
		authFormWithValidationErrors := templates.AuthForm(dto.Username, dto.Email, true, validationErrors)
		c.Logger().Debug("loginDTO contains validationErrors, returning them back")
		return authFormWithValidationErrors.Render(c.Request().Context(), c.Response().Writer)
	}

	c.SetCookie(a.jwtService.NewTokenCookie(token))
	c.Redirect(http.StatusSeeOther, "/")
	return nil
}
