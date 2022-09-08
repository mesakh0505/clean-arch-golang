package domain

import "errors"

var (
	// ErrCtxNil will throw if the context in any function need ctx is nil
	ErrCtxNil = errors.New("Context is Nil")
)
