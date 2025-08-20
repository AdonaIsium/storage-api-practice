package store

import "errors"

var (
	ErrNotFOund = errors.New("not found")
	ErrConflict = errors.New("conflict")
)
