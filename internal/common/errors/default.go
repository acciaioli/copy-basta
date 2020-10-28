package errors

import (
	"fmt"
)

type Error struct {
	description   string
	isUserError   bool
	internalError error
}

type ErrorOption func(e *Error)

func WithDescription(format string, a ...interface{}) ErrorOption {
	return func(e *Error) {
		e.description = fmt.Sprintf(format, a...)
	}
}

func WithInternalError(internalError error) ErrorOption {
	return func(e *Error) {
		e.internalError = internalError
	}
}

func IsUserError() ErrorOption {
	return func(e *Error) {
		e.isUserError = true
	}
}

func NewError(opts ...ErrorOption) *Error {
	e := Error{
		description:   "unknown error",
		isUserError:   false,
		internalError: nil,
	}
	for _, opt := range opts {
		opt(&e)
	}
	return &e
}

func (e *Error) Error() string {
	return ""
}

type ErrorBuilder interface {
	NewUserErrorf(format string, a ...interface{}) error
	NewUserError(err error) error
	NewInternalErrorf(format string, a ...interface{}) error
	NewInternalError(err error) error
}

type errorBuilder struct {
	module string
}

func NewErrorBuilder(module string) ErrorBuilder {
	return &errorBuilder{
		module: module,
	}
}

func (e errorBuilder) NewUserErrorf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	return fmt.Errorf("user error\n[%s] %w", e.module, err)
}

func (e errorBuilder) NewUserError(err error) error {
	return e.NewUserErrorf("%w", err)
}

func (e errorBuilder) NewInternalErrorf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	return fmt.Errorf("internal error\n[%s] %w", e.module, err)
}

func (e errorBuilder) NewInternalError(err error) error {
	return e.NewInternalErrorf("%w", err)
}
