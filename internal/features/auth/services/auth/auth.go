package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jim-ww/nms-go/internal/features/auth/dtos"
	"github.com/jim-ww/nms-go/internal/features/auth/role"
	"github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	validtr "github.com/jim-ww/nms-go/internal/features/auth/validator"
	"github.com/jim-ww/nms-go/internal/repository"
)

const (
	usernameF string = "username"
	passwordF string = "password"
	emailF    string = "email"
)

const (
	usernameTaken        = "username already exists"
	emailTaken           = "email already exists"
	usernameDoesNotExist = "username does not exist"
	invalidPassword      = "invalid password"
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
	fieldErrors := convertValidationErrorsToMap(vErrors)

	if vErrors != nil {
		return "", fieldErrors, nil
	}

	user, err := srv.repo.FindUserByUsername(ctx, dto.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fieldErrors[usernameF] = append(fieldErrors[usernameF], usernameDoesNotExist)
			return "", fieldErrors, nil
		}
		return "", nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	if err := ComparePasswords(user.Password, dto.Password); err != nil {
		fieldErrors[passwordF] = append(fieldErrors[passwordF], invalidPassword)
		return "", fieldErrors, nil
	}

	jwtToken, err = srv.jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return jwtToken, nil, nil
}

func (srv *AuthService) RegisterUser(ctx context.Context, dto *dtos.RegisterDTO) (jwtToken string, validationErrors map[string][]string, err error) {

	vErrors := validtr.ValidateRegisterDTO(dto)
	fieldErrors := convertValidationErrorsToMap(vErrors)

	if vErrors != nil {
		return "", fieldErrors, nil
	}

	taken, err := srv.repo.IsUsernameTaken(ctx, dto.Username)
	if err != nil {
		return "", nil, err
	} else if taken == 1 {
		fieldErrors[usernameF] = append(fieldErrors[usernameF], usernameTaken)
	}

	taken, err = srv.repo.IsEmailTaken(ctx, dto.Email)
	if err != nil {
		return "", nil, err
	} else if taken == 1 {
		fieldErrors[emailF] = append(fieldErrors[emailF], emailTaken)
	}

	fmt.Println("fieldErrors len:", len(fieldErrors)) // TODO remove

	if len(fieldErrors) > 0 { // TODO remove len(fieldErrors[usernameF]) > 0 || len(fieldErrors[emailF]) > 0 || len(fieldErrors[passwordF]) > 0 {
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

func convertValidationErrorsToMap(validationErrors validator.ValidationErrors) map[string][]string {
	vErrs := map[string][]string{}
	for _, vErr := range validationErrors {
		vErrs[vErr.Field()] = append(vErrs[vErr.Field()], vErr.Error())
	}
	return vErrs
}
