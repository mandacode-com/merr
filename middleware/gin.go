package merrmid

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
)

// ErrorResponse represents the JSON structure for error responses
type ErrorResponse struct {
	Error string         `json:"error"`
	Code  merr.ErrCode   `json:"code"`
}

// GinErrorHandler is a Gin middleware that handles errors and converts them to JSON responses.
// It processes all errors in the context and returns the first merr.PublicErr found,
// or a generic internal server error if no public errors are found.
func GinErrorHandler() gin.HandlerFunc {
	return GinErrorHandlerWithOptions(nil)
}

// GinErrorHandlerOptions provides configuration options for the error handler
type GinErrorHandlerOptions struct {
	// LogErrors determines whether to log internal errors (default: true)
	LogErrors bool
	// CustomErrorResponse allows customizing the error response format
	CustomErrorResponse func(c *gin.Context, publicErr merr.PublicErr)
	// OnInternalError is called when a non-public error occurs
	OnInternalError func(c *gin.Context, err error)
}

// GinErrorHandlerWithOptions creates a Gin error handler with custom options
func GinErrorHandlerWithOptions(opts *GinErrorHandlerOptions) gin.HandlerFunc {
	if opts == nil {
		opts = &GinErrorHandlerOptions{
			LogErrors: true,
		}
	}

	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		// Process errors and find the first public error
		var publicErr merr.PublicErr
		var internalErr error

		for _, ginErr := range c.Errors {
			if pe, ok := ginErr.Err.(merr.PublicErr); ok {
				if publicErr == nil {
					publicErr = pe
				}
			} else if internalErr == nil {
				internalErr = ginErr.Err
			}
		}

		// If we found a public error, use it
		if publicErr != nil {
			if opts.CustomErrorResponse != nil {
				opts.CustomErrorResponse(c, publicErr)
			} else {
				c.JSON(publicErr.Code().ToHTTPStatus(), ErrorResponse{
					Error: publicErr.Public(),
					Code:  publicErr.Code(),
				})
			}
			return
		}

		// Handle internal error
		if internalErr != nil {
			if opts.LogErrors {
				log.Printf("Internal error: %v", internalErr)
			}
			
			if opts.OnInternalError != nil {
				opts.OnInternalError(c, internalErr)
			} else {
				c.JSON(merr.ErrInternalServerError.ToHTTPStatus(), ErrorResponse{
					Error: "Internal server error",
					Code:  merr.ErrInternalServerError,
				})
			}
		}
	}
}

// AbortWithError is a helper function to abort with a merr.PublicErr
func AbortWithError(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}

// AbortWithPublicError is a helper function to create and abort with a public error
func AbortWithPublicError(c *gin.Context, code merr.ErrCode, public string, baseErr error) {
	err := merr.New(code, public, baseErr)
	AbortWithError(c, err)
}
