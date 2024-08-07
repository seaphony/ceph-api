package types

import (
	"errors"
)

var (
	ErrNotImplemented  = errors.New("NotImplemented")
	ErrInvalidConfig   = errors.New("InvalidConfig")
	ErrInvalidArg      = errors.New("InvalidArg")
	ErrNotFound        = errors.New("NotFound")
	ErrAlreadyExists   = errors.New("ErrAlreadyExists")
	ErrInternal        = errors.New("InternalError")
	ErrUnauthenticated = errors.New("Unauthenticated")
	ErrAccessDenied    = errors.New("AccessDenied")
)
