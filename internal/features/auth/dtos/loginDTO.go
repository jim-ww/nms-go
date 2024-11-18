package dtos

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
