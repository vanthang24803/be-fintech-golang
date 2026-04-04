package apperr

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrConflict            = errors.New("resource already exists")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidInput        = errors.New("invalid input")
	ErrInternal            = errors.New("internal server error")
)
