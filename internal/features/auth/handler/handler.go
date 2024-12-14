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
	dto := dtos.NewLoginDTO(req.FormValue("username"), req.FormValue("password"))

	token, validationErrors, err := a.authService.LoginUser(c.Request().Context(), dto)
	if err != nil {
		return fmt.Errorf("unable to login user %w", err)
	}

	fmt.Printf("validationErrors(map) len is %d", len(validationErrors)) // TODO remove

	if len(validationErrors) > 0 {
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
	dto := dtos.NewRegisterDTO(req.FormValue("username"), req.FormValue("email"), req.FormValue("password"))

	token, validationErrors, err := a.authService.RegisterUser(c.Request().Context(), dto)
	if err != nil {
		return fmt.Errorf("unable to register user: %w", err)
	}

	if len(validationErrors) > 0 {
		authFormWithValidationErrors := templates.AuthForm(dto.Username, dto.Email, true, validationErrors)
		return authFormWithValidationErrors.Render(c.Request().Context(), c.Response().Writer)
	}

	c.SetCookie(a.jwtService.NewTokenCookie(token))
	c.Redirect(http.StatusSeeOther, "/")
	return nil
}
