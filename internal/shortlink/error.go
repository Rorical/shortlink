package shortlink

import "errors"

var (
	ErrAlreadyExists     = errors.New("Record AlreadyExists")
	ErrDoesNotExists     = errors.New("Record NotExists")
	ErrIllegalParameters = errors.New("Illegal Parameters")
)
