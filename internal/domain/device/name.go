package device

import (
	"errors"
	"regexp"
	"strings"
)

type Name struct {
	value string
}

var nameRegex = regexp.MustCompile(`^[a-zA-Z0-9 _.-]+$`)

func NewName(value string) (Name, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return Name{}, errors.New("device name is required")
	}

	if len(value) < 3 {
		return Name{}, errors.New("device name must be at least 3 characters")
	}
	if len(value) > 50 {
		return Name{}, errors.New("device name must be at most 50 characters")
	}

	if !nameRegex.MatchString(value) {
		return Name{}, errors.New("device name may only contain letters, numbers, underscores, periods, or hyphens")
	}

	return Name{value: value}, nil
}

func (n Name) String() string {
	return n.value
}
