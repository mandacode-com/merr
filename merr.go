// Package merr
// This package provides a custom error type for handling errors in a structured way.
package merr

type PublicErr interface {
	error
	Code() ErrCode
	Public() string
	Unwrap() error
}

type Err struct {
	code      ErrCode
	publicMsg string
	message   string
	cause     error
}

// Unwrap implements PublicErr.
func (e *Err) Unwrap() error {
	return e.cause
}

// Code implements PublicErr.
func (e *Err) Code() ErrCode {
	return e.code
}

// Error implements PublicErr.
func (e *Err) Error() string {
	return e.message
}

// Public implements PublicErr.
func (e *Err) Public() string {
	return e.publicMsg
}

// newPublic creates a new public error with the given code, public message, detailed message, and cause.
func newPublic(code ErrCode, publicMsg string, message string, error error) PublicErr {
	return &Err{
		code:      code,
		publicMsg: publicMsg,
		message:   message,
		cause:     error,
	}
}

// New creates a new error
func New(code ErrCode, publicMsg string, message string, error error) error {
	return newPublic(code, publicMsg, message, error)
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	if serr, ok := err.(*Err); ok {
		return &Err{
			code:      serr.code,
			publicMsg: serr.publicMsg,
			message:   message,
			cause:     serr,
		}
	}
	return &Err{
		code:      ErrUnknown,
		publicMsg: "An unexpected error occurred",
		message:   message,
		cause:     err,
	}
}

func Trace(err error) []string {
	if err == nil {
		return nil
	}
	var trace []string
	for {
		trace = append(trace, err.Error())
		if serr, ok := err.(*Err); ok && serr.cause != nil {
			err = serr.cause
		} else {
			break
		}
	}
	return trace
}

// MapPublicErr maps a PublicErr to a standard error type.
func MapPublicErr(err error) PublicErr {
	if err == nil {
		return nil
	}
	if serr, ok := err.(PublicErr); ok {
		return serr
	}
	return &Err{
		code:      ErrUnknown,
		publicMsg: "An unexpected error occurred",
		message:   err.Error(),
		cause:     err,
	}
}

// IsPublicErr checks if an error is a PublicErr.
func IsPublicErr(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(PublicErr)
	return ok
}
