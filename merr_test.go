// merr_test.go
package merr

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Local test codes
const (
	codeRoot ErrCode = "E_ROOT"
)

//go:noinline
func rootErr() error {
	baseErr := errors.New("root failed")
	return New(codeRoot, "root error", baseErr)
}

func TestNew_CreatesPublicErr(t *testing.T) {
	err := rootErr()

	pe, ok := err.(PublicErr)
	require.True(t, ok, "New should return a PublicErr")

	assert.Equal(t, "root error", pe.Public(), "Public() should return the message")
	assert.Equal(t, codeRoot, pe.Code(), "Code() should return the correct code")
}

func TestCheckCode(t *testing.T) {
	err := rootErr()
	
	assert.True(t, CheckCode(err, codeRoot), "CheckCode should return true for matching code")
	assert.False(t, CheckCode(err, ErrNotFound), "CheckCode should return false for non-matching code")
	assert.False(t, CheckCode(nil, codeRoot), "CheckCode should return false for nil error")
}

func TestErrorString(t *testing.T) {
	baseErr := errors.New("root failed")
	err := New(codeRoot, "root error", baseErr)
	
	assert.Equal(t, "root failed", err.Error(), "Error() should return the underlying error message")
	
	pe, ok := err.(PublicErr)
	require.True(t, ok, "Should be PublicErr")
	assert.Equal(t, "root error", pe.Public(), "Public() should return the public message")
}

func TestHTTPStatusMapping(t *testing.T) {
	assert.Equal(t, 404, ErrNotFound.ToHTTPStatus(), "ErrNotFound should map to 404")
	assert.Equal(t, 400, ErrBadRequest.ToHTTPStatus(), "ErrBadRequest should map to 400")
	assert.Equal(t, 500, ErrInternalServerError.ToHTTPStatus(), "ErrInternalServerError should map to 500")
}

func TestGRPCCodeMapping(t *testing.T) {
	assert.Equal(t, 5, int(ErrNotFound.ToGRPCCode()), "ErrNotFound should map to NotFound gRPC code")
	assert.Equal(t, 3, int(ErrInvalidInput.ToGRPCCode()), "ErrInvalidInput should map to InvalidArgument gRPC code")
}
