package fpcnv

import "fmt"

type Errorer interface {
	Error() string
	Unwrap() error
	SetWrapped(err error) Errorer
}

type Error struct {
	message string
	wrapped error
}

func NewError(message string) Errorer {
	return &Error{message, nil}
}

func NewErrorf(format string, args ...any) Errorer {
	return NewError(fmt.Sprintf(format, args...))
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) Unwrap() error {
	return e.wrapped
}

func (e *Error) SetWrapped(err error) Errorer {
	e.wrapped = err
	return e
}
