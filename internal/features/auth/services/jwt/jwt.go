package jwt

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jim-ww/nms-go/internal/config"
	"github.com/jim-ww/nms-go/internal/features/auth/role"
	"github.com/jim-ww/nms-go/internal/utils/loggers/sl"
)

var (
	JWTTokenCookieName = "jwt-token"
	ErrInvalidJWT      = errors.New("failed to validate JWT")
)

type AuthClaims struct {
	UserID uuid.UUID
	Role   role.Role
	jwt.RegisteredClaims
}

type JWTService struct {
	logger *slog.Logger
	cfg    *config.JWTTokenConfig
}

func New(logger *slog.Logger, cfg *config.JWTTokenConfig) *JWTService {
	return &JWTService{
		logger: logger,
		cfg:    cfg,
	}
}

func (srv JWTService) GenerateToken(userID uuid.UUID, role role.Role) (encodedToken string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, AuthClaims{UserID: userID, Role: role, RegisteredClaims: jwt.RegisteredClaims{
		Issuer:    "nms",
		Subject:   "user-auth",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(srv.cfg.ExpirationDuration)),
	}})
	return token.SignedString(srv.cfg.Secret)
}

func (srv JWTService) NewTokenCookie(jwtToken string) *http.Cookie {
	return &http.Cookie{
		Name:     JWTTokenCookieName,
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}
}

func (srv JWTService) ValidateAndExtractPayload(encodedToken string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(encodedToken, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return srv.cfg.Secret, nil
	})
	if err != nil {
		srv.logger.Error("Failed to parse jwt", sl.Err(err))
	} else if claims, ok := token.Claims.(*AuthClaims); ok {
		srv.logger.Debug("parsed jwt", slog.String("User ID", claims.UserID.String()), slog.String("Role", string(claims.Role)))
		return claims, nil
	}
	return nil, errors.New("unknown claims type, cannot proceed")
}
