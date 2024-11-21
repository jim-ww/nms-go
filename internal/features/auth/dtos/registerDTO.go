package dtos

import "log/slog"

type RegisterDTO struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewRegisterDTO(username, email, password string) *RegisterDTO {
	return &RegisterDTO{
		Username: username,
		Email:    email,
		Password: password,
	}
}

func (dto *RegisterDTO) SlogAttr() slog.Attr {
	return slog.Group("register_dto",
		slog.String("username", dto.Username),
		slog.String("email", dto.Email),
	)
}
