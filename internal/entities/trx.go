package entities

import "context"

type Trx interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}
