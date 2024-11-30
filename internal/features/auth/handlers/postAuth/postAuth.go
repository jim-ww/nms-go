package postauth

import (
	"log/slog"
	"net/http"
	"text/template"

	"github.com/jim-ww/nms-go/internal/features/auth/dtos"
	"github.com/jim-ww/nms-go/internal/features/auth/services/auth"
	"github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	"github.com/jim-ww/nms-go/internal/templates"
	"github.com/jim-ww/nms-go/internal/utils/loggers/sl"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *auth.AuthService
	jwtService  *jwt.JWTService
	logger      *slog.Logger
	tmpl        *template.Template
}

func New(log *slog.Logger, userService *auth.AuthService, jwtService *jwt.JWTService) *AuthHandler {
	return &AuthHandler{
		logger:      log,
		authService: userService,
		jwtService:  jwtService,
	}
}

func (ah *AuthHandler) Login(c echo.Context) error {
	// parse form & convert form data to dto
	req := c.Request()
	if err := req.ParseForm(); err != nil {
		ah.logger.Error("unable to parse form", sl.Err(err))
		return err
	}
	dto := dtos.NewLoginDTO(req.FormValue("username"), req.FormValue("password"))

	// attempt to login
	token, validationErrors, err := ah.authService.LoginUser(c.Request().Context(), dto)
	ah.logger.Error("unable to login user", sl.Err(err))
	if err != nil {
		return err
	}

	// if has errors, return validation errors to form
	if validationErrors.HasErrors() {
		ah.logger.Error("login dto has validation errors", sl.Err(err))
		return templates.AuthForm(dto.Username, "", false, validationErrors.TranslateToMap()).Render(c.Request().Context(), c.Response().Writer)
	}

	// if no errors, set token cookie and redirect to home page
	c.SetCookie(ah.jwtService.NewTokenCookie(token))
	c.Redirect(http.StatusSeeOther, "/")
	return nil
}

func (lh *AuthHandler) Register(c echo.Context) error {
	// parse form & convert form data to dto
	req := c.Request()
	if err := req.ParseForm(); err != nil {
		lh.logger.Error("unable to parse form", sl.Err(err))
		return err
	}
	dto := dtos.NewRegisterDTO(req.FormValue("username"), req.FormValue("email"), req.FormValue("password"))

	// attempt to login
	token, validationErrors, err := lh.authService.RegisterUser(c.Request().Context(), dto)
	lh.logger.Error("unable to register user", sl.Err(err))
	if err != nil {
		return err
	}

	// if has errors, return validation errors to form
	if validationErrors.HasErrors() {
		lh.logger.Error("register dto has validation errors", sl.Err(err))
		return templates.AuthForm(dto.Username, dto.Email, true, validationErrors.TranslateToMap()).Render(c.Request().Context(), c.Response().Writer)
	}

	// if no errors, set token cookie and redirect to home page
	c.SetCookie(lh.jwtService.NewTokenCookie(token))
	c.Redirect(http.StatusSeeOther, "/")
	return nil
}
