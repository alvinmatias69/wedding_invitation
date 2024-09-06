package entities

import "time"

type Token struct {
	Id         uint64    `db:"id"`
	JwtToken   string    `db:"jwt_token"`
	SteamToken string    `db:"steam_token"`
	ClaimedAt  time.Time `db:"claimed_at"`
}
