package user

import (
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
