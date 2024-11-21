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

	sessionData, ok := data["token"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("token data missing or invalid")
	}

	if exp, ok := sessionData["ExpirationTime"].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339, exp)
		if err != nil {
			return nil, fmt.Errorf("invalid expiration time: %w", err)
		}
		payload.ExpirationTime = parsedTime
	} else {
		return nil, fmt.Errorf("expiration time missing or invalid")
	}

	if issuedAt, ok := sessionData["IssuedAt"].(float64); ok {
		payload.IssuedAt = int64(issuedAt)
	} else {
		return nil, fmt.Errorf("issued at missing or invalid")
	}

	if subject, ok := sessionData["Subject"].(string); ok {
		payload.Subject = subject
	} else {
		return nil, fmt.Errorf("subject missing or invalid")
	}

	if userID, ok := sessionData["UserId"].(float64); ok {
		payload.UserId = int64(userID)
	} else {
		return nil, fmt.Errorf("user ID missing or invalid")
	}

	if role, ok := sessionData["Role"].(string); ok {
		switch role {
		case string(user.ROLE_USER), string(user.ROLE_ADMIN):
			payload.Role = user.Role(role)
		default:
			return nil, fmt.Errorf("invalid role: %s", role)
		}
	} else {
		return nil, fmt.Errorf("role missing or invalid")
	}

	return payload, nil
}
