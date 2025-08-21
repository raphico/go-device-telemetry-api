package user

import (
	"errors"
	"regexp"
	"strings"
)

type Username struct {
	value string
}

func NewUsername(value string) (Username, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return Username{}, errors.New("username is required")
	}

	if len(value) < 3 {
		return Username{}, errors.New("username must be at least 3 characters")
	}

	if len(value) > 30 {
		return Username{}, errors.New("username must be at most 30 characters")
	}

	valid := regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	if !valid.MatchString(value) {
		return Username{}, errors.New("username may only contain letters, numbers, _, ., and -")
	}

	return Username{value: value}, nil
}

func (u Username) String() string {
	return u.value
}
