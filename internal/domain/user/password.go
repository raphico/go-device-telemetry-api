package user

import (
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
		return ErrPasswordRequired
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
		return ErrPasswordTooShort
	}
	if !hasLetter || !hasDigit || !hasSpecial {
		return ErrPasswordTooWeak
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
