package merrmid

import (
	"context"

	"github.com/mandacode-com/merr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// GRPCError is a gRPC middleware that handles errors and converts them to gRPC status errors.
func GRPCError() grpc.UnaryServerInterceptor {
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

		if publicErr, ok := err.(merr.PublicErr); ok {
			return nil, status.Errorf(
				publicErr.Code().ToGRPCCode(),
				"%s", publicErr.Public(),
			)
		}
		return nil, status.Errorf(
			status.Code(err),
			"%s", err.Error(),
		)
	}
}
