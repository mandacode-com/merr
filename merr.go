// Package merr
// This package provides a custom error type for handling errors in a structured way.
package merr

// PublicErr is an interface for errors that can be publicly displayed.
type PublicErr interface {
	error
	Public() string
	Code() ErrCode
}

type err struct {
	error
	public string
	code   ErrCode
}

// New creates a new error with a public message.
func New(code ErrCode, public string, error error) error {
	return &err{
		error:  error,
		public: public,
		code:   code,
	}
}

// Public returns the public message of the error.
func (e *err) Public() string {
	return e.public
}

// Code returns the error code.
func (e *err) Code() ErrCode {
	return e.code
}
