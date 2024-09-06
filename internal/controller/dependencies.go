package controller

import (
	"context"
	"io"

	"github.com/alvinmatias69/wedding_invitation/internal/entities"
)

type jwtResource interface {
	GenerateToken(context.Context, map[string]interface{}) (string, error)
	ParseToken(context.Context, string) (entities.JwtPayload, error)
}

type exifResource interface {
	Embed(context.Context, map[string]interface{}) (func(io.Writer) error, error)
}

type tokenRepository interface {
	GetByJwtToken(context.Context, string) (entities.Token, error)
	BeginTrx(context.Context) (trx, error)
	FindOneUnclaimed(context.Context, trx) (entities.Token, error)
	Claim(context.Context, trx, uint64, string) error
}

type trx interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}
