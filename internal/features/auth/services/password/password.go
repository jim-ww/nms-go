package password

type PasswordHasher interface {
	HashPassword(password string) (hashedPassword string, err error)
	ComparePasswords(hashedPassword, password string) (err error)
}
