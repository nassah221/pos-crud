package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secret string) (Maker, error) {
	if len(secret) < minSecretKeySize {
		return nil, fmt.Errorf("secret must be at least %d characters long", minSecretKeySize)
	}

	return &JWTMaker{secret}, nil
}

func (j *JWTMaker) CreateToken(username string, cashierID int32, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, cashierID, duration)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrTokenInvalid
		}
		return []byte(j.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		vErr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(vErr.Inner, ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrTokenInvalid
	}

	return payload, nil
}
