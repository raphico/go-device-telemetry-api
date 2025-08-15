package user

import "errors"

var (
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrUsernameAlreadyExits = errors.New("username already exists")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrInvalidUsername      = errors.New("invalid username")
	ErrInvalidPassword      = errors.New("invalid password")
)
