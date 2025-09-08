package merrmid

import (
	"context"
	"errors"
	"testing"

	"github.com/mandacode-com/merr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCErrorInterceptor_PublicError(t *testing.T) {
	interceptor := GRPCErrorInterceptor()
	
	handler := func(ctx context.Context, req any) (any, error) {
		baseErr := errors.New("database error")
		return nil, merr.New(merr.ErrNotFound, "Resource not found", baseErr)
	}
	
	resp, err := interceptor(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"},
		handler,
	)
	
	assert.Nil(t, resp)
	require.Error(t, err)
	
	st, ok := status.FromError(err)
	require.True(t, ok)
	
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, "Resource not found", st.Message())
}

func TestGRPCErrorInterceptor_InternalError(t *testing.T) {
	interceptor := GRPCErrorInterceptor()
	
	handler := func(ctx context.Context, req any) (any, error) {
		return nil, errors.New("internal error")
	}
	
	resp, err := interceptor(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"},
		handler,
	)
	
	assert.Nil(t, resp)
	require.Error(t, err)
	
	st, ok := status.FromError(err)
	require.True(t, ok)
	
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, "Internal server error", st.Message())
}

func TestGRPCErrorInterceptor_Success(t *testing.T) {
	interceptor := GRPCErrorInterceptor()
	
	handler := func(ctx context.Context, req any) (any, error) {
		return "success", nil
	}
	
	resp, err := interceptor(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"},
		handler,
	)
	
	assert.Equal(t, "success", resp)
	assert.NoError(t, err)
}

func TestGRPCErrorInterceptor_ExistingStatusError(t *testing.T) {
	interceptor := GRPCErrorInterceptor()
	
	existingErr := status.Errorf(codes.Aborted, "operation aborted")
	handler := func(ctx context.Context, req any) (any, error) {
		return nil, existingErr
	}
	
	resp, err := interceptor(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"},
		handler,
	)
	
	assert.Nil(t, resp)
	assert.Equal(t, existingErr, err)
}

func TestNewPublicError(t *testing.T) {
	baseErr := errors.New("base error")
	err := NewPublicError(merr.ErrBadRequest, "Public message", baseErr)
	
	publicErr, ok := err.(merr.PublicErr)
	require.True(t, ok)
	
	assert.Equal(t, "Public message", publicErr.Public())
	assert.Equal(t, merr.ErrBadRequest, publicErr.Code())
	assert.Equal(t, "base error", err.Error())
}