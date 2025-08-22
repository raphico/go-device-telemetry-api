package user

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)
)

// ---------- Types ----------

type UserID uuid.UUID

type Email struct {
	value string
}

type Username struct {
	value string
}

type Password struct {
	hash []byte
}

type User struct {
	ID        UserID
	Email     Email
	Username  Username
	Password  Password
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ---------- UserID ----------

func NewUserID(id string) (UserID, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return UserID(uuid.Nil), err
	}

	return UserID(parsed), nil
}

func (u UserID) String() string {
	return uuid.UUID(u).String()
}

// ---------- Email ----------

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

// ---------- Username ----------

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

	if !usernameRegex.MatchString(value) {
		return Username{}, errors.New("username may only contain letters, numbers, _, ., and -")
	}

	return Username{value: value}, nil
}

func (u Username) String() string {
	return u.value
}

// ---------- Password ----------

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

// ---------- User ----------

func NewUser(email Email, username Username, password Password) *User {
	return &User{
		Email:    email,
		Username: username,
		Password: password,
	}
}

// ---------- Rehydration ----------

func RehydrateUser(
	id uuid.UUID,
	emailStr string,
	usernameStr string,
	passwordHash []byte,
	createdAt, updatedAt time.Time,
) (*User, error) {
	e, err := NewEmail(emailStr)
	if err != nil {
		return nil, fmt.Errorf("corrupt email: %w", err)
	}

	uname, err := NewUsername(usernameStr)
	if err != nil {
		return nil, fmt.Errorf("corrupt username: %w", err)
	}

	return &User{
		ID:        UserID(id),
		Email:     e,
		Username:  uname,
		Password:  PasswordFromHash(passwordHash),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}
