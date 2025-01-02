package errors

import "errors"

var (
	// Route access
	ErrUnauthorized = errors.New("unauthorized to access route")

	// User Service
	ErrUserAlreadyExists    = errors.New("username or email already exists")
	ErrUsernameDoesNotExist = errors.New("username does not exist")

	// JWT
	ErrInvalidJWT    = errors.New("failed to validate JWT")
	ErrTokenExpired  = errors.New("token has expired")
	ErrUnknownClaims = errors.New("unknown claims type, cannot proceed")
)
