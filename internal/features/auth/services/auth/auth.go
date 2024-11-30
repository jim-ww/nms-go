package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jim-ww/nms-go/internal/features/auth"
	"github.com/jim-ww/nms-go/internal/features/auth/dtos"
	"github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	"github.com/jim-ww/nms-go/internal/features/auth/services/password"
	"github.com/jim-ww/nms-go/internal/features/user"
	"github.com/jim-ww/nms-go/internal/features/user/storage"
	"github.com/jim-ww/nms-go/internal/repository"
	"github.com/jim-ww/nms-go/internal/utils/loggers/sl"
)

// TODO use context
// TODO make all SQL related stuff in (if possible readonly) transactions
type AuthService struct {
	logger    *slog.Logger
	jwt       *jwt.JWTService
	pwdHasher password.PasswordHasher
	userRepo  *repository.Queries
	validatr  *auth.AuthValidator
}

func New(logger *slog.Logger, jwtService *jwt.JWTService, passwordHasher password.PasswordHasher, userRepo *repository.Queries, validatr *auth.AuthValidator) *AuthService {
	return &AuthService{
		logger:    logger,
		jwt:       jwtService,
		pwdHasher: passwordHasher,
		userRepo:  userRepo,
		validatr:  validatr,
	}
}

func (srv *AuthService) LoginUser(ctx context.Context, dto *dtos.LoginDTO) (jwtToken string, validationErrors auth.ValidationErrors, err error) {

	// validate dto
	if validationErrors = srv.validatr.ValidateLoginDTO(dto); validationErrors.HasErrors() {
		return "", validationErrors, nil
	}

	user, err := srv.userRepo.FindUserByUsername(ctx, dto.Username)
	if err != nil {

		if errors.Is(err, storage.ErrUsernameDoesNotExist) {
			srv.logger.Debug("Failed to get user by username", sl.Err(err))
			validationErrors[auth.UsernameField] = append(validationErrors[auth.UsernameField], auth.UsernameDoesNotExist)
			return "", validationErrors, nil
		}

		srv.logger.Error("Failed to get user by username", sl.Err(err))
		return "", nil, err
	}

	if err := srv.pwdHasher.ComparePasswords(user.Password, dto.Password); err != nil {
		srv.logger.Debug("Password hash comparison failure", sl.Err(err))
		validationErrors[auth.PasswordField] = append(validationErrors[auth.PasswordField], auth.InvalidPassword)
		return "", validationErrors, nil
	}

	srv.logger.Debug("Generating JWT token")
	jwtToken, err = srv.jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		srv.logger.Error("Failed to generate JWT token", sl.Err(err))
		return "", nil, err
	}

	return jwtToken, nil, nil
}

func (srv *AuthService) RegisterUser(ctx context.Context, dto *dtos.RegisterDTO) (jwtToken string, validationErrors auth.ValidationErrors, err error) {

	if validationErrors = srv.validatr.ValidateRegisterDTO(dto); validationErrors.HasErrors() {
		srv.logger.Debug("field validation completed with errors:", slog.Any("validationErrors", validationErrors))
		return "", validationErrors, nil
	}

	srv.logger.Debug("field validation completed, checking for existing username/email")

	taken, err := srv.userRepo.IsUsernameTaken(ctx, dto.Username)
	if err != nil {
		return "", nil, err
	} else if taken {
		validationErrors[auth.UsernameField] = append(validationErrors[auth.UsernameField], auth.UsernameTaken)
	}
	srv.logger.Debug("username validated")

	taken, err = srv.userRepo.IsEmailTaken(dto.Email)
	if err != nil {
		return "", nil, err
	} else if taken {
		validationErrors[auth.EmailField] = append(validationErrors[auth.EmailField], ErrEmailTaken)
	}

	srv.logger.Debug("email validated")

	if validationErrors.HasErrors() {
		srv.logger.Debug("field validation completed with errors:", slog.Any("validationErrors", validationErrors))
		return "", validationErrors, nil
	}

	srv.logger.Debug("field validation completed")

	srv.logger.Debug("Generating hashed password")
	hashedPassword, err := srv.pwdHasher.HashPassword(dto.Password)
	if err != nil {
		srv.logger.Error("Failed to generate hash for password", sl.Err(err))
		return "", nil, err
	}
	dto.Password = string(hashedPassword)
	srv.logger.Debug("User attemt to register:", dto.SlogAttr())

	srv.logger.Debug("Creating user with user repository")
	userID, err := srv.userRepo.Create(dto.Username, dto.Email, string(hashedPassword), user.ROLE_USER)
	if err != nil {
		srv.logger.Error("Failed to create new user", sl.Err(err))
		return "", nil, err
	}

	srv.logger.Debug("Generating JWT token")
	jwtToken, err = srv.jwt.GenerateToken(userID, user.ROLE_USER)
	if err != nil {
		srv.logger.Error("Failed to generate JWT token", sl.Err(err))
		return "", nil, err
	}

	return jwtToken, nil, nil
}
