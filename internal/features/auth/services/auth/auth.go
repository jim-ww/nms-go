package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jim-ww/nms-go/internal/features/auth/dtos"
	"github.com/jim-ww/nms-go/internal/features/auth/role"
	"github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	validtr "github.com/jim-ww/nms-go/internal/features/auth/validator"
	"github.com/jim-ww/nms-go/internal/repository"
)

const (
	usernameF string = "Username"
	passwordF string = "Password"
	emailF    string = "Email"
)

const (
	usernameTaken        = "Username already exists"
	emailTaken           = "Email already exists"
	usernameDoesNotExist = "Username does not exist"
	invalidPassword      = "Invalid password"
)

var (
	ErrUserAlreadyExists    = errors.New("username or email already exists")
	ErrUsernameDoesNotExist = errors.New("username does not exist")
)

// TODO use context
// TODO make all SQL related stuff in (if possible readonly) transactions
type AuthService struct {
	jwt  *jwt.JWTService
	repo *repository.Queries
}

func New(jwtService *jwt.JWTService, repo *repository.Queries) *AuthService {
	return &AuthService{
		jwt:  jwtService,
		repo: repo,
	}
}

func (srv *AuthService) LoginUser(ctx context.Context, dto *dtos.LoginDTO) (jwtToken string, validationErrors map[string][]string, err error) {
	vErrors := validtr.ValidateLoginDTO(dto)
	validationErrors = validtr.ConvertValidationErrorsToMap(vErrors)

	if len(validationErrors) > 0 {
		return "", validationErrors, nil
	}

	user, err := srv.repo.FindUserByUsername(ctx, dto.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			validationErrors[usernameF] = append(validationErrors[usernameF], usernameDoesNotExist)
			return "", validationErrors, nil
		}
		return "", nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	if err = ComparePasswords(user.Password, dto.Password); err != nil {
		validationErrors[passwordF] = append(validationErrors[passwordF], invalidPassword)
		return "", validationErrors, nil
	}

	jwtToken, err = srv.jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return jwtToken, nil, nil
}

func (srv *AuthService) RegisterUser(ctx context.Context, dto *dtos.RegisterDTO) (jwtToken string, validationErrors map[string][]string, err error) {

	vErrors := validtr.ValidateRegisterDTO(dto)
	validationErrors = validtr.ConvertValidationErrorsToMap(vErrors)

	if len(validationErrors) > 0 {
		return "", validationErrors, nil
	}

	taken, err := srv.repo.IsUsernameTaken(ctx, dto.Username)
	if err != nil {
		return "", nil, err
	} else if taken == 1 {
		validationErrors[usernameF] = append(validationErrors[usernameF], usernameTaken)
	}

	taken, err = srv.repo.IsEmailTaken(ctx, dto.Email)
	if err != nil {
		return "", nil, err
	} else if taken == 1 {
		validationErrors[emailF] = append(validationErrors[emailF], emailTaken)
	}

	if len(validationErrors) > 0 {
		return "", validationErrors, nil
	}

	hashedPassword, err := HashPassword(dto.Password)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate hash for password: %w", err)
	}

	userID, err := srv.repo.InsertUser(ctx, repository.InsertUserParams{
		ID:        uuid.New(),
		Username:  dto.Username,
		Email:     dto.Email,
		Password:  string(hashedPassword),
		Role:      role.ROLE_USER,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return "", nil, fmt.Errorf("failed to create new user: %w", err)
	}

	jwtToken, err = srv.jwt.GenerateToken(userID.ID, role.ROLE_USER)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return jwtToken, nil, nil
}
