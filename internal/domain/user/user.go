package user

import (
	"fmt"
	"time"
)

type User struct {
	ID        UserID
	Email     Email
	Username  Username
	Password  Password
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email Email, username Username, password Password) *User {
	return &User{
		Email:    email,
		Username: username,
		Password: password,
	}
}

func RehydrateUser(
	id UserID,
	emailStr string,
	usernameStr string,
	passwordHash []byte,
	createdAt, updatedAt time.Time,
) (*User, error) {
	e, err := NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("corrupt email: %w", err)
	}
	uname, err := NewUsername(usernameStr)
	if err != nil {
		return nil, fmt.Errorf("corrupt username: %w", err)
	}

	return &User{
		ID:        id,
		Email:     e,
		Username:  uname,
		Password:  PasswordFromHash(passwordHash),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
