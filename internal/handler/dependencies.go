package handler

import (
	"context"
	"io"
)

type controller interface {
	GetHiddenImage(context.Context, io.Writer) error
}
