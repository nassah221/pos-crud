package token

import "time"

type Maker interface {
	CreateToken(username string, id int32, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
