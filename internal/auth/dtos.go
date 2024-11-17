package auth

type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password" validate:"email"`
}

type RegisterDTO struct {
	Username string `json:"username" validate:"min=3 max=30"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=3 max=255"`
}
