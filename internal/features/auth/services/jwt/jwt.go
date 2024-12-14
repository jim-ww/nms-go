package jwt

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jim-ww/nms-go/internal/config"
	"github.com/jim-ww/nms-go/internal/features/auth/role"
)

var (
	JWTTokenCookieName = "jwt-token"
	ErrInvalidJWT      = errors.New("failed to validate JWT")
	ErrTokenExpired    = errors.New("token has expired")
	ErrUnknownClaims   = errors.New("unknown claims type, cannot proceed")
)

type AuthClaims struct {
	UserID uuid.UUID
	Role   role.Role
	jwt.RegisteredClaims
}

type JWTService struct {
	cfg *config.JWTTokenConfig
}

func New(cfg *config.JWTTokenConfig) *JWTService {
	return &JWTService{
		cfg: cfg,
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
		return nil, err
	} else if claims, ok := token.Claims.(*AuthClaims); ok {
		return claims, nil
	}
	return nil, ErrUnknownClaims
}
