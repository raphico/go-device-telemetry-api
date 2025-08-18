package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        UserID
	Email     Email
	Username  Username
	Password  Password
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, username, password string) (*User, error) {
	addr, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	uname, err := NewUsername(username)
	if err != nil {
		return nil, err
	}

	pass, err := NewPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		Email:    addr,
		Username: uname,
		Password: pass,
	}, nil
}

func RehydrateUser(
	id uuid.UUID,
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
		ID:        UserID(id),
		Email:     e,
		Username:  uname,
		Password:  PasswordFromHash(passwordHash),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
