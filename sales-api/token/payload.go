package token

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrTokenExpired = fmt.Errorf("token has expired")
var ErrTokenInvalid = fmt.Errorf("token is invalid")

type Payload struct {
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	Username  string    `json:"username"`
	CashierID int32     `json:"cashierID"`
	ID        uuid.UUID `json:"id"`
}

func NewPayload(username string, cashierID int32, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		CashierID: cashierID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrTokenExpired
	}
	return nil
}
