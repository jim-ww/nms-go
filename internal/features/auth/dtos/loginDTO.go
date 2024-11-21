package dtos

import "log/slog"

type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewLoginDTO(username, password string) *LoginDTO {
	return &LoginDTO{
		Username: username,
		Password: password,
	}
}

func (dto *LoginDTO) SlogAttr() slog.Attr {
	return slog.Group("login_dto",
		slog.String("username", dto.Username),
	)
}
