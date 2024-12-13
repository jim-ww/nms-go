package dtos

type LoginDTO struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=3,max=255"`
}

type RegisterDTO struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=255"`
}

func NewLoginDTO(username, password string) *LoginDTO {
	return &LoginDTO{Username: username, Password: password}
}

func NewRegisterDTO(username, email, password string) *RegisterDTO {
	return &RegisterDTO{Username: username, Email: email, Password: password}
}
