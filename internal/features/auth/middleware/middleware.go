package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jim-ww/nms-go/internal/features/auth/role"
	"github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	"github.com/labstack/echo/v4"
)

var RouteAccessLevels = map[string]AccessLevel{
	"/error":           AllowUnauthorized,
	"/favicon.ico":     AllowUnauthorized,
	"/login":           OnlyUnauthorized,
	"/register":        OnlyUnauthorized,
	"/api/login":       OnlyUnauthorized,
	"/api/register":    OnlyUnauthorized,
	"/":                OnlyAuthorized,
	"/api/logout":      OnlyAuthorized,
	"/api/notes":       OnlyAuthorized,
	"/api/user":        OnlyAuthorized,
	"/admin":           OnlyAdmins,
	"/api/admin/users": OnlyAdmins,
	"/api/admin/notes": OnlyAdmins,
}

var DynamicRouteAccessLevels = map[string]AccessLevel{
	"/web/":       AllowUnauthorized,
	"/api/notes/": OnlyAuthorized,
	"/api/admin/": OnlyAdmins,
}

var (
	ErrAccessLevelNotFound = errors.New("failed to determine route access level")
	ErrUnauthorized        = errors.New("unauthorized")
)

type AuthMiddleware struct {
	jwtService *jwt.JWTService
}

func New(jwtService *jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (a AuthMiddleware) Handler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jwtCookie, err := c.Cookie(jwt.JWTTokenCookieName)
		if err != nil {
			c.Logger().Debug("No JWT cookie found, err:", err)
			return a.handleUnauthorized(next)(c)
		}

		token, err := a.jwtService.ValidateAndExtractPayload(jwtCookie.Value)
		if err != nil {
			c.Logger().Errorf("Failed to extract JWT payload", err)
			return err
		}

		// TODO validate somewhere else?
		if time.Now().After(token.ExpiresAt.Time) {
			c.Logger().Info("JWT has expired")
			return jwt.ErrTokenExpired
		}

		if err = checkIsAllowedToAccess(c, c.Request().URL.Path, token.Role); err != nil {
			return err
		}

		return next(c)
	}
}

func checkIsAllowedToAccess(c echo.Context, urlPath string, role role.Role) error {
	accessLevel, isRouteAcessLevelFound := RouteAccessLevels[urlPath]
	if !isRouteAcessLevelFound {
		isDynamicRouteAcessLevelFound := false
		accessLevel, isDynamicRouteAcessLevelFound = CheckDynamicPath(urlPath)
		if !isDynamicRouteAcessLevelFound {
			c.Logger().Error("failed to determine route access level (no url path defined in routeAccessLevel/dynamicAccessLevel maps)")
			return ErrAccessLevelNotFound
		}
	}
	if !IsAllowed(role, accessLevel) {
		return ErrUnauthorized
	}
	return nil
}

func (am AuthMiddleware) handleUnauthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		urlPath := c.Request().URL.Path
		accessLevel, ok := RouteAccessLevels[urlPath]
		if !ok {
			accessLevel, ok = CheckDynamicPath(urlPath)
			if !ok {
				return ErrAccessLevelNotFound
			}
		}

		if accessLevel == AllowUnauthorized || accessLevel == OnlyUnauthorized {
			next(c)
			return nil
		}

		c.Redirect(http.StatusSeeOther, "/login")
		return nil
	}
}

func CheckDynamicPath(pathToCheck string) (accessLvl AccessLevel, ok bool) {
	for path, access := range DynamicRouteAccessLevels {
		if strings.HasPrefix(pathToCheck, path) {
			return access, true
		}
	}
	return accessLvl, false
}
func IsAllowed(r role.Role, access AccessLevel) bool {
	switch access {
	case AllowUnauthorized:
		return true
	case OnlyUnauthorized:
		return false
	case OnlyAuthorized:
		return r == role.ROLE_USER || r == role.ROLE_ADMIN
	case OnlyAdmins:
		return r == role.ROLE_ADMIN
	default:
		return false
	}
}
