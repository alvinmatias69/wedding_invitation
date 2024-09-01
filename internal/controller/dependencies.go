package controller

import (
	"context"
	"io"
)

type jwtService interface {
	GetToken(context.Context) (string, error)
}

type exifService interface {
	EmbedAndWrite(context.Context, string, io.Writer) error
}
