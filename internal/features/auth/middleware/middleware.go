package middleware

import (
	"net/http"

	"github.com/jim-ww/nms-go/internal/features/auth/role"
	"github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	"github.com/jim-ww/nms-go/internal/lib/errors"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	jwtService *jwt.JWTService
}

func New(jwtService *jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (a *AuthMiddleware) OnlyAdmins(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jwtCookie, err := c.Cookie(jwt.JWTCookieName)
		if err != nil {
			c.Logger().Debug("JWT cookie is not present")
			return errors.ErrUnauthorized
		}

		payload, err := a.jwtService.ValidateAndExtractPayload(jwtCookie.Value)
		if err != nil {
			c.Logger().Debug("failed to extract jwt payload, err:", err)
			c.Redirect(http.StatusSeeOther, "/logout")
			return nil
		}

		if payload.Role == role.ROLE_ADMIN {
			return next(c)
		}

		return errors.ErrUnauthorized
	}
}

func (a *AuthMiddleware) OnlyUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jwtCookie, err := c.Cookie(jwt.JWTCookieName)
		if err != nil {
			c.Logger().Debug("JWT cookie is not present")
			return errors.ErrUnauthorized
		}

		payload, err := a.jwtService.ValidateAndExtractPayload(jwtCookie.Value)
		if err != nil {
			c.Logger().Debug("failed to extract jwt payload, err:", err)
			c.Redirect(http.StatusSeeOther, "/logout")
			return nil
		}

		if payload.Role == role.ROLE_USER || payload.Role == role.ROLE_ADMIN {
			return next(c)
		}

		return errors.ErrUnauthorized
	}
}

func (a *AuthMiddleware) OnlyUnauthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := c.Cookie(jwt.JWTCookieName)
		if err != nil {
			return next(c)
		}
		// if user has JWT cookie, then redirect him
		return c.Redirect(http.StatusSeeOther, "/")
	}
}
