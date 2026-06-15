package apperrors

import "errors"

var (
	ErrTodoNotFound      = errors.New("todo not found")
	ErrTaskTitleEmpty    = errors.New("task cannot be empty")
	ErrDueDateInvalid    = errors.New("invalid due_date")
	ErrTodoIDNotPositive = errors.New("invalid id")
)
