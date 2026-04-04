package apperr

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
)
