package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID    pgtype.UUID `json:"id"`
	Email string      `json:"email"`

	jwt.RegisteredClaims
}

func NewPayload(id pgtype.UUID, email string, duration time.Duration) (*Payload, error) {

	now := time.Now()

	payload := &Payload{
		ID:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}

	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiresAt.Time) {
		return ErrExpiredToken
	}

	return nil
}
