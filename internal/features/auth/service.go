package auth

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jim-ww/nms-go/internal/features/auth/dtos"
	"github.com/jim-ww/nms-go/internal/features/user"
	"github.com/jim-ww/nms-go/pkg/config"
	"github.com/jim-ww/nms-go/pkg/utils/jwts"
	"golang.org/x/crypto/bcrypt"
)

// TODO add password hashing
// TODO implement more robust logging
// TODO use context?
// TODO make all SQL related stuff in readonly(if possible) transactions

const (
	ErrUsernameTaken = "username already exists"
	ErrEmailTaken    = "email already exists"
)

var (
	ErrInvalidJWT = errors.New("failed to validate JWT")
)

type AuthRepository interface {
	IsUsernameTaken(username string) (taken bool, err error)
	IsEmailTaken(email string) (taken bool, err error)
	CreateUser(username, email, hashedPassword string, role user.Role) (createdID int64, err error)
}

type AuthService struct {
	logger *slog.Logger
	cfg    *config.JWTTokenConfig
	repo   AuthRepository
}

func NewAuthService(logger *slog.Logger, cfg *config.JWTTokenConfig, repo AuthRepository) *AuthService {
	return &AuthService{
		logger: logger,
		cfg:    cfg,
		repo:   repo,
	}
}

type token struct {
	ExpirationTime int64 // TODO how to store expiration time in token?
	IssuedAt       int64
	Subject        string
	UserId         int64
	RoleName       string // TODO use 'Role' type here?
}

func (service AuthService) NewToken(userID int64, role user.Role) (encodedToken string, err error) {
	token := token{
		ExpirationTime: service.cfg.ExpirationTime.Microseconds(),
		IssuedAt:       time.Now().Unix(),
		Subject:        "user-auth",
		UserId:         userID,
		RoleName:       string(role), // TODO
	}
	claims := map[string]any{"session": token}

	return jwts.GenerateJWT(service.cfg.Secret, claims)
}

func NewTokenCookie(jwtToken string) *http.Cookie {
	return &http.Cookie{
		Name:     "jwt-token",
		Value:    jwtToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // TODO
	}
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

func (srv *AuthService) RegisterUser(dto *dtos.RegisterDTO) (jwtToken string, validationErrors map[string][]string, err error) {
	validationErrors = make(map[string][]string)

	if validationErrors := ValidateRegisterDTO(dto); len(validationErrors) > 0 {
		srv.logger.Info("field validation completed with errors:", slog.Any("validationErrors", validationErrors))
		return "", validationErrors, nil
	}

	srv.logger.Info("field validation completed, checking for existing username/email")

	taken, err := srv.repo.IsUsernameTaken(dto.Username)
	if err != nil {
		return "", nil, err
	} else if taken {
		validationErrors[UsernameField] = append(validationErrors[UsernameField], ErrUsernameTaken)
	}

	taken, err = srv.repo.IsEmailTaken(dto.Email)
	if err != nil {
		return "", nil, err
	} else if taken {
		validationErrors[EmailField] = append(validationErrors[EmailField], ErrEmailTaken)
	}

	// TODO test
	fmt.Println("val errors:", validationErrors, "validation errors len:", len(validationErrors))
	if len(validationErrors) > 0 {
		return "", validationErrors, nil
	}

	srv.logger.Info("field validation completed")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}
	dto.Password = string(hashedPassword)
	srv.logger.Info("User attemt to register:", slog.Any("registerDTO:", dto))

	userID, err := srv.repo.CreateUser(dto.Username, dto.Email, string(hashedPassword), user.ROLE_USER)
	if err != nil {
		return "", nil, err
	}

	jwtToken, err = srv.NewToken(userID, user.ROLE_USER)
	if err != nil {
		return "", nil, err
	}

	return jwtToken, nil, nil
}
