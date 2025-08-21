package user

import (
	"errors"
	"net/mail"
	"strings"
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return Email{}, errors.New("email is required")
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return Email{}, errors.New("invalid email format")
	}

	return Email{value: value}, nil
}

func (e Email) String() string {
	return e.value
}
