package jwts

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

func ValidateAndExtractPayload(secret, token string) (map[string]any, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrNotEnoughParts
	}

	headerEncoded, claimsEncoded, signatureEncoded := parts[0], parts[1], parts[2]

	signatureInput := fmt.Sprintf("%s.%s", headerEncoded, claimsEncoded)
	h := hmac.New(sha256.New, []byte(secret))
	_, err := h.Write([]byte(signatureInput))
	if err != nil {
		return nil, err
	}
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	// validate signature
	if !hmac.Equal([]byte(expectedSignature), []byte(signatureEncoded)) {
		return nil, ErrInvalidSignature
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(claimsEncoded)
	if err != nil {
		return nil, err
	}

	var claims map[string]any
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, err
	}

	return claims, nil
}
