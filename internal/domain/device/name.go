package device

import (
	"regexp"
	"strings"
)

type Name struct {
	value string
}

func NewName(value string) (Name, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return Name{}, ErrNameRequired
	}

	if len(value) < 3 {
		return Name{}, ErrNameTooShort
	}
	if len(value) > 50 {
		return Name{}, ErrNameTooLong
	}

	valid := regexp.MustCompile(`^[a-zA-Z0-9 _.-]+$`)
	if !valid.MatchString(value) {
		return Name{}, ErrNameInvalidChars
	}

	return Name{value: value}, nil
}

func (n Name) String() string {
	return n.value
}
