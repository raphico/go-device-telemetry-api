package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUsernameTaken      = errors.New("username already exists")

	ErrUsernameRequired     = errors.New("username is required")
	ErrUsernameTooShort     = errors.New("username must be at least 3 characters")
	ErrUsernameTooLong      = errors.New("username must be at most 30 characters")
	ErrUsernameInvalidChars = errors.New("username contains invalid characters (only letters, numbers, underscores allowed)")

	ErrEmailRequired = errors.New("email is required")
	ErrEmailInvalid  = errors.New("invalid email format")

	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordTooWeak  = errors.New("password must contain a mix of letters, numbers, and symbols")
)
