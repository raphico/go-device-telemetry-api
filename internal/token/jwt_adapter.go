package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTGenerator interface {
	Generate(userID string, exp time.Duration) (string, error)
	Validate(tokenStr string) (*claims, error)
}

type JWTAdapter struct {
	secret []byte
}

type claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTAdapter(secret []byte) JWTGenerator {
	return &JWTAdapter{secret: secret}
}

func (j *JWTAdapter) Generate(userID string, exp time.Duration) (string, error) {
	c := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(j.secret)
}

func (j *JWTAdapter) Validate(tokenStr string) (*claims, error) {
	claims := claims{}

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (any, error) {
		// Ensure the signing method is HS256
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrWrongTokenType
		}

		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return &claims, nil
}
