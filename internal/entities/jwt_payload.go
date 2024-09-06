package entities

import "time"

type JwtPayload struct {
	IssuedAt time.Time
	TokenId  string
}
