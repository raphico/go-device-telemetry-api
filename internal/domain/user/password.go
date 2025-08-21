package user

import (
	"errors"
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hash []byte
}

func PasswordFromHash(hash []byte) Password {
	return Password{hash: hash}
}

func (p Password) Hash() string {
	return string(p.hash)
}

func (p *Password) Matches(candidate string) bool {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(candidate))
	return err == nil
}

func validatePassword(pw string) error {
	if len(pw) == 0 {
		return errors.New("password is required")
	}

	var (
		hasMinLen  = len(pw) >= 8
		hasLetter  bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, c := range pw {
		switch {
		case unicode.IsLetter(c):
			hasLetter = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsSymbol(c) || unicode.IsPunct(c):
			hasSpecial = true
		}
	}

	if !hasMinLen {
		return errors.New("password must be at least 8 characters")
	}
	if !hasLetter || !hasDigit || !hasSpecial {
		return errors.New("password must contain a mix of letters, numbers, and symbols")
	}

	return nil
}

func NewPassword(value string) (Password, error) {
	if err := validatePassword(value); err != nil {
		return Password{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, fmt.Errorf("failed to hash password: %w", err)
	}

	return Password{hash: hash}, nil
}
