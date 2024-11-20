package jwt

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/jim-ww/nms-go/internal/features/user"
	"github.com/jim-ww/nms-go/pkg/config"
	"github.com/jim-ww/nms-go/pkg/utils/jwts"
	"github.com/jim-ww/nms-go/pkg/utils/loggers/sl"
)

var (
	JWTTokenCookieName = "jwt-token"
	ErrInvalidJWT      = errors.New("failed to validate JWT")
)

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

func (srv JWTService) GenerateToken(userID int64, role user.Role) (encodedToken string, err error) {
	token := Payload{
		ExpirationTime: time.Now().Add(srv.cfg.ExpirationDuration),
		IssuedAt:       time.Now().Unix(),
		Subject:        "user-auth",
		UserId:         userID,
		Role:           role,
	}
	claims := map[string]any{"session": token}

	return jwts.GenerateJWT(srv.cfg.Secret, claims)
}

func (srv JWTService) NewTokenCookie(jwtToken string) *http.Cookie {
	return &http.Cookie{
		Name:     JWTTokenCookieName,
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,                    // TODO
		SameSite: http.SameSiteDefaultMode, // TODO
	}
}

func (srv JWTService) ValidateAndExtractPayload(encodedToken string) (*Payload, error) {
	srv.logger.Debug("Extraction payload from jwt...")
	payloadData, err := jwts.ValidateAndExtractPayload(srv.cfg.Secret, encodedToken)
	if err != nil {
		srv.logger.Error("Failed to extract payload data from jwt", sl.Err(err))
		return nil, ErrInvalidJWT
	}
	srv.logger.Debug("Successfully validated JWT")
	srv.logger.Debug("Mapping payload to Payload struct...")

	payload, err := MapToPayload(payloadData)
	if err != nil {
		srv.logger.Error("Failed to map payload to Payload struct", sl.Err(err))
		return nil, err
	}

	return payload, nil
}
