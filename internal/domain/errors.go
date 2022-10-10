package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrDuplicateKey    = errors.New("duplicate key")
	ErrInvalidArgument = errors.New("invalid argument")
)
