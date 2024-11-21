package bcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct{}

func New() PasswordHasher {
	return PasswordHasher{}
}

func (ps PasswordHasher) HashPassword(password string) (hashedPassword string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func (ps PasswordHasher) ComparePasswords(hashedPassword, password string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
