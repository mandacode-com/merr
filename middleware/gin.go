package merrmid

import (
	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
)

// GinErrorHandler is a Gin middleware that handles errors and converts them to JSON responses.
func GinErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				if publicErr, ok := err.Err.(merr.PublicErr); ok {
					c.JSON(publicErr.Code().ToHTTPStatus(), gin.H{
						"error": publicErr.Public(),
						"code":  publicErr.Code(),
					})
					return
				} else {
					c.JSON(merr.ErrInternalServerError.ToHTTPStatus(), gin.H{
						"error": "Internal server error",
						"code":  merr.ErrInternalServerError,
					})
				}
			}
		}
	}
}
