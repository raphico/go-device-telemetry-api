package user

import (
	"net/mail"
	"strings"
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)

	// Required
	if value == "" {
		return Email{}, ErrEmailRequired
	}

	// Validate format
	if _, err := mail.ParseAddress(value); err != nil {
		return Email{}, ErrEmailInvalid
	}

	return Email{value: value}, nil
}

func (e Email) String() string {
	return e.value
}
