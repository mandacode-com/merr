package merr

import "net/http"

// error map
var httpErrorMap = map[ErrCode]int{
	ErrUnknown:                     http.StatusInternalServerError,
	ErrNotFound:                    http.StatusNotFound,
	ErrInvalidInput:                http.StatusBadRequest,
	ErrPermissionDenied:            http.StatusForbidden,
	ErrInternalServerError:         http.StatusInternalServerError,
	ErrTimeout:                     http.StatusGatewayTimeout,
	ErrConflict:                    http.StatusConflict,
	ErrUnauthorized:                http.StatusUnauthorized,
	ErrBadRequest:                  http.StatusBadRequest,
	ErrServiceUnavailable:          http.StatusServiceUnavailable,
	ErrTooManyRequests:             http.StatusTooManyRequests,
	ErrGatewayTimeout:              http.StatusGatewayTimeout,
	ErrUnprocessableEntity:         http.StatusUnprocessableEntity,
	ErrNotImplemented:              http.StatusNotImplemented,
	ErrMethodNotAllowed:            http.StatusMethodNotAllowed,
	ErrForbidden:                   http.StatusForbidden,
	ErrPreconditionFailed:          http.StatusPreconditionFailed,
	ErrExpectationFailed:           http.StatusExpectationFailed,
	ErrBadGateway:                  http.StatusBadGateway,
	ErrLengthRequired:              http.StatusLengthRequired,
	ErrUnsupportedMediaType:        http.StatusUnsupportedMediaType,
	ErrRangeNotSatisfiable:         http.StatusRequestedRangeNotSatisfiable,
	ErrInsufficientStorage:         http.StatusInsufficientStorage,
	ErrLoopDetected:                http.StatusLoopDetected,
	ErrNotAcceptable:               http.StatusNotAcceptable,
	ErrTooEarly:                    http.StatusTooEarly,
	ErrRequestHeaderFieldsTooLarge: http.StatusRequestHeaderFieldsTooLarge,
}

func (e ErrCode) ToHTTPStatus() int {
	if status, exists := httpErrorMap[e]; exists {
		return status
	}
	return http.StatusInternalServerError // Default to Internal Server Error if the error code is not mapped
}
