package merrmid

import (
	"context"
	"log"

	"github.com/mandacode-com/merr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCErrorInterceptor is a gRPC middleware that intercepts errors and converts them to gRPC status errors.
func GRPCErrorInterceptor() grpc.UnaryServerInterceptor {
	return GRPCErrorInterceptorWithOptions(nil)
}

// GRPCErrorInterceptorOptions provides configuration options for the gRPC error interceptor
type GRPCErrorInterceptorOptions struct {
	// LogErrors determines whether to log internal errors (default: true)
	LogErrors bool
	// OnInternalError is called when a non-public error occurs
	OnInternalError func(ctx context.Context, err error) error
	// OnPublicError is called when a public error occurs, allows customization
	OnPublicError func(ctx context.Context, publicErr merr.PublicErr) error
}

// GRPCErrorInterceptorWithOptions creates a gRPC error interceptor with custom options
func GRPCErrorInterceptorWithOptions(opts *GRPCErrorInterceptorOptions) grpc.UnaryServerInterceptor {
	if opts == nil {
		opts = &GRPCErrorInterceptorOptions{
			LogErrors: true,
		}
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		// Handle merr.PublicErr
		if publicErr, ok := err.(merr.PublicErr); ok {
			if opts.OnPublicError != nil {
				if customErr := opts.OnPublicError(ctx, publicErr); customErr != nil {
					return nil, customErr
				}
			}
			
			return nil, status.Errorf(
				publicErr.Code().ToGRPCCode(),
				"%s", publicErr.Public(),
			)
		}

		// Handle other errors
		if opts.LogErrors {
			log.Printf("gRPC internal error in %s: %v", info.FullMethod, err)
		}

		if opts.OnInternalError != nil {
			if customErr := opts.OnInternalError(ctx, err); customErr != nil {
				return nil, customErr
			}
		}

		// Check if error is already a gRPC status error
		if _, ok := status.FromError(err); ok {
			return nil, err
		}

		// Convert to internal gRPC error
		return nil, status.Errorf(codes.Internal, "Internal server error")
	}
}

// GRPCStreamErrorInterceptor handles errors for streaming gRPC calls
func GRPCStreamErrorInterceptor() grpc.StreamServerInterceptor {
	return GRPCStreamErrorInterceptorWithOptions(nil)
}

// GRPCStreamErrorInterceptorWithOptions creates a streaming gRPC error interceptor with custom options
func GRPCStreamErrorInterceptorWithOptions(opts *GRPCErrorInterceptorOptions) grpc.StreamServerInterceptor {
	if opts == nil {
		opts = &GRPCErrorInterceptorOptions{
			LogErrors: true,
		}
	}

	return func(
		srv any,
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := handler(srv, stream)
		if err == nil {
			return nil
		}

		// Handle merr.PublicErr
		if publicErr, ok := err.(merr.PublicErr); ok {
			if opts.OnPublicError != nil {
				if customErr := opts.OnPublicError(stream.Context(), publicErr); customErr != nil {
					return customErr
				}
			}
			
			return status.Errorf(
				publicErr.Code().ToGRPCCode(),
				"%s", publicErr.Public(),
			)
		}

		// Handle other errors
		if opts.LogErrors {
			log.Printf("gRPC stream internal error in %s: %v", info.FullMethod, err)
		}

		if opts.OnInternalError != nil {
			if customErr := opts.OnInternalError(stream.Context(), err); customErr != nil {
				return customErr
			}
		}

		// Check if error is already a gRPC status error
		if _, ok := status.FromError(err); ok {
			return err
		}

		// Convert to internal gRPC error
		return status.Errorf(codes.Internal, "Internal server error")
	}
}

// NewPublicError is a helper function to create a new public error
func NewPublicError(code merr.ErrCode, public string, baseErr error) error {
	return merr.New(code, public, baseErr)
}
