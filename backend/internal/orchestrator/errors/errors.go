package errors

import "fmt"

var (
	ErrConflict     = fmt.Errorf("resource already exists")
	ErrBadRequest   = fmt.Errorf("bad request")
	ErrUnauthorized = fmt.Errorf("unauthorized")
)
