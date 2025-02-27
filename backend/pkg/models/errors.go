package models

import "errors"

var (
	ErrNoTasks                = errors.New("no tasks")
	ErrTaskDoesNotExist       = errors.New("task does not exist")
	ErrUnsupportedOperation   = errors.New("unsupported operation operation")
	ErrDivisionByZero         = errors.New("division by zero")
	ErrInvalidExpression      = errors.New("invalid expression")
	ErrNoExpressions          = errors.New("no expressions are found")
	ErrExpressionDoesNotExist = errors.New("expression does not exist")
)
