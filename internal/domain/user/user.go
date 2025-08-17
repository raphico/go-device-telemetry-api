package user

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserID uuid.UUID

type Email struct {
	value string
}

type Username struct {
	value string
}

type User struct {
	ID        UserID
	Email     Email
	Username  Username
	Password  Password
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Password struct {
	hash []byte
}

func newUser(email, username, password string) (*User, error) {
	addr, err := NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	uname, err := NewUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	pass, err := newPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &User{
		Email:    addr,
		Username: uname,
		Password: pass,
	}, nil
}

func NewUserID(id string) (UserID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return UserID(uuid.Nil), nil
	}

	return UserID(parsed), nil
}

func (u UserID) String() string {
	return uuid.UUID(u).String()
}

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)

	if _, err := mail.ParseAddress(value); err != nil {
		return Email{}, fmt.Errorf("%w: %s", ErrInvalidEmail, value)
	}

	return Email{value: value}, nil
}

func (e Email) String() string {
	return e.value
}

func NewUsername(value string) (Username, error) {
	value = strings.TrimSpace(value)

	if len(value) < 3 || len(value) > 30 {
		return Username{}, fmt.Errorf("%w: username must be between 3 and 30 characters", ErrInvalidUsername)
	}

	valid := regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
	if !valid.MatchString(value) {
		return Username{}, fmt.Errorf("%w: username may only contain letters, numbers, underscores, periods, and hyphens", ErrInvalidUsername)
	}

	return Username{value: value}, nil
}

func (u Username) String() string {
	return u.value
}

func newPassword(value string) (Password, error) {
	if err := validatePassword(value); err != nil {
		return Password{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, fmt.Errorf("failed to hash password: %w", err)
	}

	return Password{hash: hash}, nil
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
	var (
		hasMinLen  = len(pw) >= 8
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, c := range pw {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsSymbol(c) || unicode.IsPunct(c):
			hasSpecial = true
		}
	}

	if !hasMinLen {
		return fmt.Errorf("%w: password must be at least 8 characters long", ErrInvalidPassword)
	}
	if !hasUpper {
		return fmt.Errorf("%w: password must contain at least one uppercase letter", ErrInvalidPassword)
	}
	if !hasLower {
		return fmt.Errorf("%w: password must contain at least one lowercase letter", ErrInvalidPassword)
	}
	if !hasDigit {
		return fmt.Errorf("%w: password must contain at least one digit", ErrInvalidPassword)
	}
	if !hasSpecial {
		return fmt.Errorf("%w: password must contain at least one special character", ErrInvalidPassword)
	}

	return nil
}
