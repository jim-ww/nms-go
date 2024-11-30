package repository

import "errors"

var (
	ErrUserAlreadyExists    = errors.New("username or email already exists")
	ErrUsernameDoesNotExist = errors.New("username does not exist")
)
