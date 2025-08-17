package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/raphico/go-device-telemetry-api/internal/domain/user"
)

type TokenID uuid.UUID

type Token struct {
	ID         TokenID
	UserID     user.UserID
	Plaintext  string
	Hash       []byte
	Scope      string
	Revoked    bool
	ExpiresAt  time.Time
	LastUsedAt time.Time
}

func NewToken(userId user.UserID, ttl time.Duration, scope string) (*Token, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenGenerationFailed, err)
	}

	plaintext := base64.RawURLEncoding.EncodeToString(b)

	hash := sha256.Sum256([]byte(plaintext))

	return &Token{
		UserID:    userId,
		Hash:      hash[:],
		ExpiresAt: time.Now().Add(ttl),
		Plaintext: plaintext,
		Scope:     scope,
	}, nil
}
