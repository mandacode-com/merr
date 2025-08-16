// Package merr
// This package provides a custom error type for handling errors in a structured way.
package merr

import (
	"errors"
	"fmt"
	"runtime"
)

// PublicErr is an interface for errors that can be publicly displayed.
type PublicErr interface {
	error
	Code() ErrCode
	Public() string
	Unwrap() error
	Stack() []uintptr
	Format(s fmt.State, verb rune)
}

type err struct {
	msg    string
	public string
	code   ErrCode
	cause  error
	stack  []uintptr
}

func (e *err) Error() string {
	return e.msg
}

func (e *err) Public() string {
	if e.public != "" {
		return e.public
	}
	return e.msg
}

func (e *err) Code() ErrCode {
	return e.code
}

func (e *err) Unwrap() error {
	if e.cause != nil {
		return e.cause
	}
	return nil
}

func (e *err) Stack() []uintptr {
	if len(e.stack) == 0 {
		return nil
	}
	return e.stack
}

func (e *err) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%s ", e.msg)
			if len(e.stack) > 0 {
				fn := runtime.FuncForPC(e.stack[0])
				if fn != nil {
					file, line := fn.FileLine(e.stack[0])
					name := fn.Name()
					fmt.Fprintf(s, "[%s:%d (%s)]\n", file, line, name)
				} else {
					fmt.Fprintf(s, "[unknown function]\n")
				}
			}
			if e.cause != nil {
				fmt.Fprintf(s, "\tcaused by: %+v", e.cause)
			}
			return
		}
		fallthrough
	case 's':
		fmt.Fprint(s, e.Public())
	default:
		fmt.Fprint(s, e.msg)
	}
}

// New creates a new error with the specified code and message.
func New(code ErrCode, msg string, publicMsg string, error error) error {
	if msg == "" {
		msg = string(code)
	}

	stack := make([]uintptr, 32)
	n := runtime.Callers(2, stack)

	e := &err{
		msg:    msg,
		public: publicMsg,
		code:   code,
		cause:  error,
		stack:  stack[:n],
	}

	return e
}

// Wrap wraps an existing error with a new message and code, capturing the stack trace.
func Wrap(error error, msg string) error {
	if error == nil {
		return nil
	}

	if msg == "" {
		msg = error.Error()
	}

	stack := make([]uintptr, 32)
	n := runtime.Callers(2, stack)

	if publicErr, ok := error.(PublicErr); ok {
		return &err{
			msg:    msg,
			public: publicErr.Public(),
			code:   publicErr.Code(),
			cause:  publicErr,
			stack:  stack[:n],
		}
	}
	return &err{
		msg:    msg,
		public: msg,
		code:   ErrUnknown,
		cause:  error,
		stack:  stack[:n],
	}
}

// Is Check if the error matches the specified code.
func Is(error error, code ErrCode) bool {
	if error == nil {
		return false
	}

	if publicErr, ok := error.(PublicErr); ok {
		return publicErr.Code() == code
	}

	return false
}

// Cause retrieves the original cause of the error, if it exists.
func Cause(error error) error {
	if error == nil {
		return nil
	}

	if publicErr, ok := error.(PublicErr); ok {
		return publicErr.Unwrap()
	}

	if cause := errors.Unwrap(error); cause != nil {
		return cause
	}

	return nil
}
