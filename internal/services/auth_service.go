package services

import (
	"errors"
	"time"

	"github.com/jim-ww/nms-go/pkg/utils"
)

var (
	ErrInvalidJWT = errors.New("failed to validate JWT")
)

type Session struct {
	ExpirationTime int64
	IssuedAt       int64
	Subject        string
	UserId         int64
	RoleName       string
}

func NewSession(userID int64, roleName string, expirationTime time.Time, secret string) (encodedSession string, err error) {
	session := Session{
		ExpirationTime: expirationTime.Unix(),
		IssuedAt:       time.Now().Unix(),
		Subject:        "user-auth",
		UserId:         userID,
		RoleName:       roleName,
	}
	claims := map[string]any{"session": session}

	return utils.GenerateJWT(secret, claims)
}

// func DecodeAndVerifyJWTSession(jwtToken, secret string) (session Session, err error) {
//
// 	parts := strings.Split(jwtToken, ".")
// 	if len(parts) != 3 {
// 		return session, ErrInvalidJWT
// 	}
// 	headerPart := parts[0]
// 	claimsPart := parts[1]
// 	signaturePart := parts[2]
//
// 	return Session{}, nil
// }
