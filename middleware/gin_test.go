package merrmid

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mandacode-com/merr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGinErrorHandler_PublicError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	r.Use(GinErrorHandler())
	
	r.GET("/test", func(c *gin.Context) {
		baseErr := errors.New("database connection failed")
		publicErr := merr.New(merr.ErrNotFound, "User not found", baseErr)
		c.Error(publicErr)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "User not found", response.Error)
	assert.Equal(t, merr.ErrNotFound, response.Code)
}

func TestGinErrorHandler_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	r.Use(GinErrorHandler())
	
	r.GET("/test", func(c *gin.Context) {
		c.Error(errors.New("internal error"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	
	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "Internal server error", response.Error)
	assert.Equal(t, merr.ErrInternalServerError, response.Code)
}

func TestGinErrorHandler_NoErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	r.Use(GinErrorHandler())
	
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAbortWithPublicError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	r := gin.New()
	r.Use(GinErrorHandler())
	
	r.GET("/test", func(c *gin.Context) {
		baseErr := errors.New("validation failed")
		AbortWithPublicError(c, merr.ErrInvalidInput, "Invalid input provided", baseErr)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "Invalid input provided", response.Error)
	assert.Equal(t, merr.ErrInvalidInput, response.Code)
}