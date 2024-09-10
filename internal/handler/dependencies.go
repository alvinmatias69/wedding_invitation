package handler

import (
	"context"
	"io"

	"github.com/alvinmatias69/wedding_invitation/internal/entities"
)

type controller interface {
	GetHiddenImage(context.Context, io.Writer) error
	GetSteamToken(context.Context, string) (entities.SteamTokenResponse, error)
	GetMessages(context.Context, uint64) ([]entities.Message, error)
	PostMessage(context.Context, entities.Message) error
}
