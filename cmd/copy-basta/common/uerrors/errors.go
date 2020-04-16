package uerrors

import (
	"errors"
	"fmt"
)

type cat int

const (
	Internal cat = iota
	Input
	Template
)

func catToString(cat cat) string {
	return []string{
		"Internal Error",
		"User Input Error",
		"Template Error",
	}[cat]
}

type Error struct {
	cat     cat
	err     error
	userMsg *string
}

type ErrorOption func(*Error)

func FromErr(err error) ErrorOption {
	return func(e *Error) {
		e.err = err
	}
}

func FromString(s string) ErrorOption {
	return func(e *Error) {
		e.err = errors.New(s)
	}
}

func newError(cat cat, userMsg *string, opts ...ErrorOption) error {
	e := &Error{cat: cat, userMsg: userMsg}

	for _, o := range opts {
		o(e)
	}

	return e
}

func (e *Error) Error() string {
	catStr := catToString(e.cat)
	if e.userMsg == nil {
		return catStr
	}
	return fmt.Sprintf("%s: %s", catStr, *e.userMsg)
}

func NewInternalError(opts ...ErrorOption) error {
	return newError(Internal, nil, opts...)
}

func NewInputError(userMsg string, opts ...ErrorOption) error {
	return newError(Input, &userMsg, opts...)
}

func NewTemplateError(userMsg string, opts ...ErrorOption) error {
	return newError(Template, &userMsg, opts...)
}

func NewFlagValidationError(flag string, reason string) error {
	return fmt.Errorf(`"--%s" %s`, flag, reason)
}
