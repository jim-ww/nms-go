package services

import (
	"github.com/jim-ww/nms-go/internal/dtos"
	"github.com/jim-ww/nms-go/internal/repositories"
	"github.com/jim-ww/nms-go/internal/validators"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func (s *UserService) RegisterUser(dto *dtos.RegisterDTO) (bool, []error) {
	isValid, validationErrors := validators.ValidateRegisterDTO(dto)
	if !isValid {
		return false, validationErrors
	}

	if s.userRepo.IsUsernameTaken(dto.Username) {
		validationErrors = append(validationErrors, validators.ErrUsernameTaken)
	}
	if s.userRepo.IsEmailTaken(dto.Email) {
		validationErrors = append(validationErrors, validators.ErrEmailTaken)
	}

	if len(validationErrors) > 0 {
		return false, validationErrors
	}

	err := s.userRepo.CreateUser(dto)
	if err != nil {
		return false, []error{err}
	}

	return true, nil
}
