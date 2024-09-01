package service

import "context"

type tokenRepository interface {
	Save(context.Context, string) error
}
