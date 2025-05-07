package errors

import "errors"

var (
	ErrNoTasks                = errors.New("no tasks")
	ErrInvalidExpression      = errors.New("invalid expression")
	ErrNoExpressions          = errors.New("no expressions are found")
	ErrExpressionDoesNotExist = errors.New("expression does not exist")
	ErrUnauthorized           = errors.New("unauthorized")
	ErrBadRequest             = errors.New("bad request")
	ErrConflict               = errors.New("conflict")
	ErrUserAlreadyExists      = errors.New("this login has already been registered")
)
