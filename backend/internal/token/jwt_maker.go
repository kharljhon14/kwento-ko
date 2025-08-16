package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func (j JWTMaker) CreateToken(id pgtype.UUID, email string, duration time.Duration) (string, error) {
	payload, err := NewPayload(id, email, duration)
	if err != nil {
		return "", nil
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    payload.ID,
		"email": payload.Email,
		"exp":   payload.ExpiresAt,
		"iat":   payload.IssuedAt,
	})

	token, err := jwt.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (j JWTMaker) VerifyToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(t *jwt.Token) (any, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Payload)
	if !ok {
		return nil, fmt.Errorf("invalid payload fields")
	}

	return claims, nil

}

func NewJWTMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("secret key must be atleast 32 characters long")
	}

	return &JWTMaker{
		secretKey: secretKey,
	}, nil
}
