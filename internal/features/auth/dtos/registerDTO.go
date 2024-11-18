package dtos

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
