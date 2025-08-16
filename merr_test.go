// merr_format_test.go
package merr

import (
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Local test codes (assuming type ErrCode is defined in the package)
const (
	codeRoot ErrCode = "E_ROOT"
)

//go:noinline
func rootErr() error {
	return New(codeRoot, "root failed", "root error", nil)
}

//go:noinline
func wrapOnce(in error) error {
	return Wrap(in, "wrap failed")
}

// topFuncName returns the function name of the top frame in pcs.
func topFuncName(pcs []uintptr) string {
	if len(pcs) == 0 {
		return ""
	}
	fn := runtime.FuncForPC(pcs[0])
	if fn == nil {
		return ""
	}
	return fn.Name()
}

func TestNew_CreatesPublicErr(t *testing.T) {
	err := rootErr()

	pe, ok := err.(PublicErr)
	require.True(t, ok, "New should return a PublicErr")

	assert.Equal(t, "root error", pe.Public(), "Public() should return the message")
	assert.Equal(t, codeRoot, pe.Code(), "Code() should return the correct code")
	assert.NotNil(t, pe.Stack(), "Stack() should return a non-nil stack trace")
}

func TestWrap_CreatesPublicErrWithCause(t *testing.T) {
	base := rootErr()
	wrapped := wrapOnce(base)

	pe, ok := wrapped.(PublicErr)
	require.True(t, ok, "Wrap should return a PublicErr")

	assert.Equal(t, "root error", pe.Public(), "Public() should return the message")
	assert.Equal(t, codeRoot, pe.Code(), "Code() should return the correct code")
	assert.NotNil(t, pe.Stack(), "Stack() should return a non-nil stack trace")

	// Cause should be the original error
	cause := Cause(pe)
	require.NotNil(t, cause, "Cause should return a non-nil error")
	assert.Equal(t, base, cause, "Cause should match the original error")
}

func TestFormatPlusV_PropagatesAndShowsCallSitePerLayer(t *testing.T) {
	base := rootErr()
	wrapped := wrapOnce(base)

	// When
	out := fmt.Sprintf("%+v", wrapped)

	// Then: top line includes top message (current layer)
	assert.Contains(t, out, "wrap failed", "top layer message should appear")

	// Should contain "caused by:" exactly once for two-layer chain
	assert.Contains(t, out, "caused by:", "should include cause propagation")

	// Top layer call site should include wrapOnce
	assert.Contains(t, out, "wrapOnce", "top-layer stack should include the wrap site function")

	// The inner layer (base) should also render its own call site via propagation
	// We specifically check that after the first "caused by:" appears, the inner message and function appear.
	causeIdx := strings.Index(out, "caused by:")
	require.GreaterOrEqual(t, causeIdx, 0, "must include caused by")
	causeTail := out[causeIdx:]

	assert.Contains(t, causeTail, "root failed", "inner layer message should propagate")
	assert.Contains(t, causeTail, "rootErr", "inner layer stack should include the root call site function")

	// Sanity: %s prints only Public() (which is "root failed" at the top layer)
	assert.Equal(t, "root error", fmt.Sprintf("%s", wrapped), "%s should print Public() only")
}

func TestFormatPlusV_SingleLayer_NoCause(t *testing.T) {
	base := rootErr()

	pe, ok := base.(PublicErr)
	require.True(t, ok, "rootErr should return PublicErr")

	// When
	out := fmt.Sprintf("%+v", base)

	// Then: contains only its own message and call site, no "caused by:"
	assert.Contains(t, out, "root failed", "single-layer message should appear")
	assert.Contains(t, out, "rootErr", "single-layer stack should include the creation site function")
	assert.NotContains(t, out, "caused by:", "single-layer should not include cause")
	// %s prints Public() == "root error"
	assert.Equal(t, "root error", fmt.Sprintf("%s", base))
	// Top frame function name matches expectation
	assert.Contains(t, topFuncName(pe.Stack()), "rootErr")
}

func TestIs_CurrentLayerOnly(t *testing.T) {
	base := rootErr()
	wrapped := wrapOnce(base)

	assert.True(t, Is(wrapped, codeRoot), "Is should be true for wrapped error")
	assert.True(t, Is(base, codeRoot), "Is should be true for base error")
}

func TestCause_ReturnsDirectCause(t *testing.T) {
	base := rootErr()
	wrapped := wrapOnce(base)

	got := Cause(wrapped)
	require.NotNil(t, got, "Cause should return the direct cause")
	assert.Equal(t, base, got, "Cause(wrap) should be the wrapped error (pointer equality expected)")

	// Base has no cause
	assert.Nil(t, Cause(base))
}
