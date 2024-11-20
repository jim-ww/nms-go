package jwt

import (
	"fmt"
	"time"

	"github.com/jim-ww/nms-go/internal/features/user"
)

type Payload struct {
	ExpirationTime time.Time
	IssuedAt       int64
	Subject        string
	UserId         int64
	Role           user.Role
}

func MapToPayload(data map[string]any) (*Payload, error) {
	payload := &Payload{}

	// Parse ExpirationTime
	if exp, ok := data["ExpirationTime"].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339, exp)
		if err != nil {
			return &Payload{}, fmt.Errorf("invalid expiration time: %w", err)
		}
		payload.ExpirationTime = parsedTime
	} else {
		return &Payload{}, fmt.Errorf("expiration time missing or invalid")
	}

	// Parse IssuedAt
	if issuedAt, ok := data["IssuedAt"].(int64); ok {
		payload.IssuedAt = int64(issuedAt)
	} else {
		return &Payload{}, fmt.Errorf("issued at missing or invalid")
	}

	// Parse Subject
	if subject, ok := data["Subject"].(string); ok {
		payload.Subject = subject
	} else {
		return &Payload{}, fmt.Errorf("subject missing or invalid")
	}

	// Parse UserId
	if userID, ok := data["UserId"].(float64); ok {
		payload.UserId = int64(userID)
	} else {
		return &Payload{}, fmt.Errorf("user ID missing or invalid")
	}

	// Parse Role
	if role, ok := data["Role"].(string); ok {
		switch role {
		case string(user.ROLE_USER), string(user.ROLE_ADMIN):
			payload.Role = user.Role(role)
		default:
			return &Payload{}, fmt.Errorf("invalid role: %s", role)
		}
	} else {
		return &Payload{}, fmt.Errorf("role missing or invalid")
	}

	return payload, nil
}
