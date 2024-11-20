package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	"github.com/jim-ww/nms-go/internal/features/user"
	"github.com/jim-ww/nms-go/pkg/utils/handlers"
)

type AccessLevel int

const (
	AllowUnauthorized AccessLevel = iota
	OnlyAuthorized
	OnlyAdmins
)

var AccessLevelString = map[AccessLevel]string{
	AllowUnauthorized: "AllowUnauthorized",
	OnlyAuthorized:    "OnlyAuthorized",
	OnlyAdmins:        "OnlyAdmins",
}

var RouteAccessLevels = map[string]AccessLevel{
	"/login":                             AllowUnauthorized,
	"/register":                          AllowUnauthorized,
	"/api/login":                         AllowUnauthorized,
	"/api/register":                      AllowUnauthorized,
	"/web/static/htmx-2.0.3/htmx.min.js": AllowUnauthorized,
	"/web/static/tailwind.css":           AllowUnauthorized,
	"/favicon.ico":                       AllowUnauthorized,
	"/":                                  OnlyAuthorized,
	"/api/user":                          OnlyAuthorized,
	"/api/logout":                        OnlyAuthorized,
	"/admin":                             OnlyAdmins,
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
			return
		}

		if time.Now().After(token.ExpirationTime) {
			am.logger.Debug("jwt token has expired, returning 401")
			handlers.RenderError(w, r, "Unauthorized", http.StatusUnauthorized)
			return
		}

		am.logger.Debug("Checking user role access")
		accessLevel, ok := RouteAccessLevels[r.URL.Path]
		if !ok {
			am.logger.Error("failed to determine route access level (no url path defined in routeAccessLevel map), returning 401", slog.String("route", r.URL.Path), slog.Any("route access", accessLevel))
			handlers.RenderError(w, r, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if IsAllowed(token.Role, accessLevel) {
			am.logger.Debug("route is allowed, returning 200", slog.String("route", r.URL.Path), slog.String("route-access", AccessLevelString[accessLevel]))
			next.ServeHTTP(w, r)
			return
		}

		am.logger.Debug("route is not allowed, returning 401", slog.String("route", r.URL.Path), slog.String("route-access", AccessLevelString[accessLevel]))
		handlers.RenderError(w, r, "Unauthorized", http.StatusUnauthorized)
	})
}

func (am AuthMiddleware) handleUnauthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessLevel, ok := RouteAccessLevels[r.URL.Path]
		if !ok {
			am.logger.Error("failed to determine route access level (no url path defined in routeAccessLevel map), returning 401", slog.String("route", r.URL.Path), slog.String("route-access", AccessLevelString[accessLevel]))
			handlers.RenderError(w, r, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if accessLevel == AllowUnauthorized {
			am.logger.Debug("route is AllowUnauthorized, returning 200", slog.String("route", r.URL.Path), slog.String("route-access", AccessLevelString[accessLevel]))
			next.ServeHTTP(w, r)
			return
		}

		am.logger.Debug("is not authorized, redirecting to /login", slog.String("route", r.URL.Path), slog.String("route-access", AccessLevelString[accessLevel]))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}

func IsAllowed(role user.Role, access AccessLevel) bool {
	switch role {
	case user.ROLE_ADMIN:
		return true
	case user.ROLE_USER:
		switch access {
		case AllowUnauthorized:
			return true
		case OnlyAuthorized:
			return true
		case OnlyAdmins:
			return false
		default:
			return false
		}
	default:
		return false
	}
}
