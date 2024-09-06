package constant

import "errors"

var (
	ErrNotFound = errors.New("Not found")
	ErrTokenExp = errors.New("Token expired")
)
