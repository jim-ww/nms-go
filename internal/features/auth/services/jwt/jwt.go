package jwt

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jim-ww/nms-go/internal/config"
	"github.com/jim-ww/nms-go/internal/features/auth/role"
	"github.com/jim-ww/nms-go/internal/lib/errors"
)

var JWTCookieName = "jwt-token"

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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims{UserID: userID, Role: role, RegisteredClaims: jwt.RegisteredClaims{
		Issuer:    "nms",
		Subject:   "user-auth",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(srv.cfg.ExpirationDuration)),
	}})
	return token.SignedString([]byte(srv.cfg.Secret))
}

func (srv JWTService) NewTokenCookie(jwtToken string) *http.Cookie {
	return &http.Cookie{
		Name:     JWTCookieName,
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		Expires:  time.Now().Add(srv.cfg.ExpirationDuration),
		SameSite: http.SameSiteLaxMode,
	}
}

func (srv JWTService) ValidateAndExtractPayload(encodedToken string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(encodedToken, &AuthClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(srv.cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.ErrInvalidJWT
	}
	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, errors.ErrUnknownClaims
	}
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.ErrTokenExpired
	}
	return claims, nil
}
