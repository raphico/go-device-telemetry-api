package user

import (
	"regexp"
	"strings"
)

type Username struct {
	value string
}

func NewUsername(value string) (Username, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return Username{}, ErrUsernameRequired
	}

	if len(value) < 3 {
		return Username{}, ErrUsernameTooShort
	}
	if len(value) > 30 {
		return Username{}, ErrUsernameTooLong
	}

	valid := regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	if !valid.MatchString(value) {
		return Username{}, ErrUsernameInvalidChars
	}

	return Username{value: value}, nil
}

func (u Username) String() string {
	return u.value
}
