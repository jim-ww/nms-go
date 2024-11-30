package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	"github.com/jim-ww/nms-go/internal/features/user"
	"github.com/jim-ww/nms-go/internal/utils/handlers"
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
	"/static/":    AllowUnauthorized,
	"/api/notes/": OnlyAuthorized,
	"/api/admin/": OnlyAdmins,
}

type AuthMiddleware struct {
	logger     *slog.Logger
	jwtService *jwt.JWTService
}

func New(logger *slog.Logger, jwtService *jwt.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		logger:     logger,
		jwtService: jwtService,
	}
}

func (am AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		am.logger.Debug("Getting jwt token from cookie")
		jwtCookie, err := r.Cookie(jwt.JWTTokenCookieName)
		if err != nil {
			am.logger.Debug("missing jwt cookie, handling as unauthorized")
			am.handleUnauthorized(next).ServeHTTP(w, r)
			return
		}

		am.logger.Debug("Validating jwt token")
		token, err := am.jwtService.ValidateAndExtractPayload(jwtCookie.Value)
		if err != nil {
			am.logger.Info("invalid jwt token, returning 401")
			handlers.RenderError(w, r, "Unauthorized", http.StatusUnauthorized)
			// TODO remove jwt cookie if invalid (after rendering error)?
			return
		}

		if time.Now().After(token.ExpiresAt.Time) {
			am.logger.Debug("jwt token has expired, returning 401")
			handlers.RenderError(w, r, "Unauthorized", http.StatusUnauthorized)
			return
		}

		am.logger.Debug("Checking needed route AccessLevel")
		accessLevel, isRouteAcessLevelFound := RouteAccessLevels[r.URL.Path]

		if !isRouteAcessLevelFound {

			am.logger.Debug("Static AccessLevel for route not found, searching dynamic paths", slog.String("route", r.URL.Path))
			isDynamicRouteAcessLevelFound := false

			accessLevel, isDynamicRouteAcessLevelFound = CheckDynamicPath(r.URL.Path)
			if !isDynamicRouteAcessLevelFound {

				am.logger.Error("failed to determine route access level (no url path defined in routeAccessLevel/dynamicAccessLevel maps), returning 401", accessLevel.SlogAttr(r.URL.Path))

				handlers.RenderError(w, r, "Unauthorized", http.StatusUnauthorized)

				return
			} else {
				am.logger.Debug("Found route AccessLevel for dynamic path", accessLevel.SlogAttr(r.URL.Path))
			}
		}

		if IsAllowed(token.Role, accessLevel) {
			am.logger.Debug("route is allowed, returning 200", slog.Any("role", token.Role), accessLevel.SlogAttr(r.URL.Path))
			next.ServeHTTP(w, r)
			return
		} else if accessLevel == OnlyUnauthorized {
			am.logger.Debug("route is allowed only for Unauthorized users, redirecting...", slog.Any("role", token.Role), accessLevel.SlogAttr(r.URL.Path))
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		am.logger.Debug("route is not allowed, returning 401", accessLevel.SlogAttr(r.URL.Path))
		handlers.RenderError(w, r, "Unauthorized", http.StatusUnauthorized)
	})
}

func (am AuthMiddleware) handleUnauthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessLevel, ok := RouteAccessLevels[r.URL.Path]
		if !ok {
			am.logger.Error("failed to determine route access level (no url path defined in routeAccessLevel map), returning 401", accessLevel.SlogAttr(r.URL.Path))
			handlers.RenderError(w, r, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if accessLevel == AllowUnauthorized || accessLevel == OnlyUnauthorized {
			am.logger.Debug("route is AllowUnauthorized, returning 200", accessLevel.SlogAttr(r.URL.Path))
			next.ServeHTTP(w, r)
			return
		}

		am.logger.Debug("is not authorized, redirecting to /login", accessLevel.SlogAttr(r.URL.Path))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}

func CheckDynamicPath(pathToCheck string) (accessLvl AccessLevel, ok bool) {
	for path, access := range DynamicRouteAccessLevels {
		if strings.HasPrefix(pathToCheck, path) {
			return access, true
		}
	}
	return accessLvl, false
}
func IsAllowed(role user.Role, access AccessLevel) bool {
	switch access {
	case AllowUnauthorized:
		return true
	case OnlyUnauthorized:
		return false
	case OnlyAuthorized:
		return role == user.ROLE_USER || role == user.ROLE_ADMIN
	case OnlyAdmins:
		return role == user.ROLE_ADMIN
	default:
		return false
	}
}
