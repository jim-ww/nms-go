package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
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

	return generateJWT(secret, claims)
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

func generateJWT(secret string, claims map[string]any) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)

	signatureInput := fmt.Sprintf("%s.%s", headerEncoded, claimsEncoded)
	h := hmac.New(sha256.New, []byte(secret))
	_, err = h.Write([]byte(signatureInput))
	if err != nil {
		return "", err
	}
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	jwt := fmt.Sprintf("%s.%s.%s", headerEncoded, claimsEncoded, signature)

	return jwt, nil
}

func verifyJWT(secret string, token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}
	headerPart := parts[0]
	claimsPart := parts[1]
	signaturePart := parts[2]

	signatureInput := fmt.Sprintf("%s.%s", headerPart, claimsPart)
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(signatureInput))
	if err != nil {
		return false
	}
	computedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return signaturePart == computedSignature
}
