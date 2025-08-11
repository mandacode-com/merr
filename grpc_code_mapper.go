package serr

import "google.golang.org/grpc/codes"

// error map
var grpcErrorMap = map[ErrCode]codes.Code{
	ErrUnknown:                     codes.Unknown,
	ErrNotFound:                    codes.NotFound,
	ErrInvalidInput:                codes.InvalidArgument,
	ErrPermissionDenied:            codes.PermissionDenied,
	ErrInternalServerError:         codes.Internal,
	ErrTimeout:                     codes.DeadlineExceeded,
	ErrConflict:                    codes.Aborted,
	ErrUnauthorized:                codes.Unauthenticated,
	ErrBadRequest:                  codes.InvalidArgument,
	ErrServiceUnavailable:          codes.Unavailable,
	ErrTooManyRequests:             codes.ResourceExhausted,
	ErrGatewayTimeout:              codes.DeadlineExceeded,
	ErrUnprocessableEntity:         codes.FailedPrecondition,
	ErrNotImplemented:              codes.Unimplemented,
	ErrMethodNotAllowed:            codes.Unimplemented,
	ErrForbidden:                   codes.PermissionDenied,
	ErrPreconditionFailed:          codes.FailedPrecondition,
	ErrExpectationFailed:           codes.FailedPrecondition,
	ErrBadGateway:                  codes.Unavailable,
	ErrLengthRequired:              codes.InvalidArgument,
	ErrUnsupportedMediaType:        codes.InvalidArgument,
	ErrRangeNotSatisfiable:         codes.OutOfRange,
	ErrInsufficientStorage:         codes.ResourceExhausted,
	ErrLoopDetected:                codes.Aborted,
	ErrNotAcceptable:               codes.InvalidArgument,
	ErrTooEarly:                    codes.FailedPrecondition,
	ErrRequestHeaderFieldsTooLarge: codes.ResourceExhausted,
}

func (e ErrCode) ToGRPCCode() codes.Code {
	if code, exists := grpcErrorMap[e]; exists {
		return code
	}
	return codes.Unknown // Default to Unknown if the error code is not mapped
}
