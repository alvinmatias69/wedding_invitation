package controller

import (
	"context"
	"io"

	"github.com/alvinmatias69/wedding_invitation/internal/entities"
	"github.com/jackc/pgx/v5"
)

type jwtResource interface {
	GenerateToken(context.Context, entities.JwtPayload) (string, error)
	ParseToken(context.Context, string) (entities.JwtPayload, error)
}

type exifResource interface {
	Embed(context.Context, map[string]interface{}) (func(io.Writer) error, error)
}

type tokenRepository interface {
	GetByJwtToken(context.Context, string) (entities.Token, error)
	BeginTrx(context.Context) (pgx.Tx, error)
	FindOneUnclaimed(context.Context, pgx.Tx) (entities.Token, error)
	Claim(context.Context, pgx.Tx, uint64, string) error
}

type messageRepository interface {
	Insert(context.Context, entities.Message) error
	Get(context.Context, uint64, uint64) ([]entities.Message, error)
}
