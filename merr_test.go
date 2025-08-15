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
	codeWrap ErrCode = "E_WRAP"
)

// Prevent inlining so that function names show up stably in stack traces.
//go:noinline
func rootErr() error {
	return New(codeRoot, "root failed", nil)
}

//go:noinline
func wrapOnce(in error) error {
	return Wrap(in, codeWrap, "wrap failed")
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

func TestFormatPlusV_PropagatesAndShowsCallSitePerLayer(t *testing.T) {
	// Given: base -> wrap
	base := rootErr()
	w := wrapOnce(base)

	// When
	out := fmt.Sprintf("%+v", w)

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
	assert.Equal(t, "root failed", fmt.Sprintf("%s", w), "%s should print Public() only")
}

func TestFormatPlusV_SingleLayer_NoCause(t *testing.T) {
	// Given: single-layer error via New
	base := rootErr()
	pe, ok := base.(PublicErr)
	require.True(t, ok, "rootErr should return PublicErr")

	// When
	out := fmt.Sprintf("%+v", base)

	// Then: contains only its own message and call site, no "caused by:"
	assert.Contains(t, out, "root failed", "single-layer message should appear")
	assert.Contains(t, out, "rootErr", "single-layer stack should include the creation site function")
	assert.NotContains(t, out, "caused by:", "single-layer should not include cause")
	// %s prints Public() == "root failed"
	assert.Equal(t, "root failed", fmt.Sprintf("%s", base))
	// Top frame function name matches expectation
	assert.Contains(t, topFuncName(pe.Stack()), "rootErr")
}

func TestIs_CurrentLayerOnly(t *testing.T) {
	base := rootErr()           // top code E_ROOT
	w := wrapOnce(base)         // top code E_WRAP

	assert.True(t, Is(w, codeWrap), "Is should be true for top code")
	assert.False(t, Is(w, codeRoot), "Is should be false for inner code on top error")

	assert.True(t, Is(base, codeRoot), "Is should match on single-layer error")
}

func TestCause_ReturnsDirectCause(t *testing.T) {
	base := rootErr()
	w := wrapOnce(base)

	got := Cause(w)
	require.NotNil(t, got, "Cause should return the direct cause")
	assert.Equal(t, base, got, "Cause(wrap) should be the wrapped error (pointer equality expected)")

	// Base has no cause
	assert.Nil(t, Cause(base))
}

func TestSetPublicMsg_OnPublicErr_ReplacesPublic_AndKeepsStack(t *testing.T) {
	// Given: a two-layer error
	base := rootErr()
	w := wrapOnce(base)

	// When: replace public message on the top-level PublicErr
	repl := SetPublicMsg(w, "custom public")
	pe, ok := repl.(PublicErr)
	require.True(t, ok, "SetPublicMsg should return a PublicErr")

	// Then: Public() replaced, internal msg kept
	assert.Equal(t, "custom public", pe.Public(), "public message should be replaced")
	assert.Equal(t, w.(PublicErr).Error(), pe.Error(), "internal message should be preserved")

	// Code should be preserved
	assert.Equal(t, w.(PublicErr).Code(), pe.Code(), "code should be preserved")

	// Cause should be preserved (direct cause of the original top-level)
	assert.Equal(t, w.(PublicErr).Unwrap(), pe.Unwrap(), "cause should be preserved")

	// Stack should be the same as the original top-level stack (not newly captured)
	origTop := w.(PublicErr).Stack()
	newTop := pe.Stack()
	require.NotEmpty(t, origTop, "original top stack should not be empty")
	require.NotEmpty(t, newTop, "new top stack should not be empty")

	origFn := topFuncName(origTop)
	newFn := topFuncName(newTop)
	assert.Equal(t, origFn, newFn, "top frame function should remain identical after SetPublicMsg")
}
