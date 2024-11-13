package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

var (
	ErrNotEnoughParts   = fmt.Errorf("not enough parts")
	ErrInvalidSignature = fmt.Errorf("invalid signature")
)

func GenerateJWT(secret string, claims map[string]any) (string, error) {
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

func VerifyJWT(secret string, token string) (claims map[string]any, isValid bool, err error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, false, ErrNotEnoughParts
	}
	headerPart := parts[0]
	claimsPart := parts[1]
	signaturePart := parts[2]

	signatureInput := fmt.Sprintf("%s.%s", headerPart, claimsPart)
	h := hmac.New(sha256.New, []byte(secret))
	if _, err = h.Write([]byte(signatureInput)); err != nil {
		return nil, false, err
	}

	claims = make(map[string]any)
	b, err := base64.RawURLEncoding.DecodeString(claimsPart)
	if err != nil {
		return nil, false, fmt.Errorf("failed to decode claims part: %w", err)
	}

	if err = json.Unmarshal(b, &claims); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal claims to JSON: %w", err)
	}

	computedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	if signaturePart == computedSignature {
		return nil, false, ErrInvalidSignature
	}

	return claims, true, nil
}
