package pagination

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	DefaultLimit = 2
	MaxLimit     = 10
)

type Cursor struct {
	ID        uuid.UUID
	CreatedAt time.Time
}

func Encode(c Cursor) string {
	payload := fmt.Sprintf("%d|%s", c.CreatedAt.UTC().UnixNano(), c.ID.String())
	return base64.RawURLEncoding.EncodeToString([]byte(payload))
}

func Decode(raw string) (Cursor, error) {
	data, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return Cursor{}, fmt.Errorf("invalid cursor encoding")
	}

	parts := strings.SplitN(string(data), "|", 2)
	if len(parts) != 2 {
		return Cursor{}, fmt.Errorf("invalid cursor format")
	}

	id, err := uuid.Parse(parts[1])
	if err != nil {
		return Cursor{}, fmt.Errorf("invalid cursor uuid")
	}

	ts, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return Cursor{}, err
	}

	return Cursor{
		CreatedAt: time.Unix(0, ts).UTC(),
		ID:        id,
	}, nil
}

func NewCursor(id uuid.UUID, createdAt time.Time) *Cursor {
	return &Cursor{
		ID:        id,
		CreatedAt: createdAt,
	}
}

func ClampLimit(n int) int {
	if n <= 0 {
		return DefaultLimit
	}
	if n > MaxLimit {
		return MaxLimit
	}
	return n
}
