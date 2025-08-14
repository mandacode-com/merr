package merr

type ErrCode string

const (
	ErrUnknown                     ErrCode = "unknown"
	ErrNotFound                    ErrCode = "not_found"
	ErrInvalidInput                ErrCode = "invalid_input"
	ErrPermissionDenied            ErrCode = "permission_denied"
	ErrInternalServerError         ErrCode = "internal_server_error"
	ErrTimeout                     ErrCode = "timeout"
	ErrConflict                    ErrCode = "conflict"
	ErrUnauthorized                ErrCode = "unauthorized"
	ErrBadRequest                  ErrCode = "bad_request"
	ErrServiceUnavailable          ErrCode = "service_unavailable"
	ErrTooManyRequests             ErrCode = "too_many_requests"
	ErrGatewayTimeout              ErrCode = "gateway_timeout"
	ErrUnprocessableEntity         ErrCode = "unprocessable_entity"
	ErrNotImplemented              ErrCode = "not_implemented"
	ErrMethodNotAllowed            ErrCode = "method_not_allowed"
	ErrForbidden                   ErrCode = "forbidden"
	ErrPreconditionFailed          ErrCode = "precondition_failed"
	ErrExpectationFailed           ErrCode = "expectation_failed"
	ErrBadGateway                  ErrCode = "bad_gateway"
	ErrLengthRequired              ErrCode = "length_required"
	ErrUnsupportedMediaType        ErrCode = "unsupported_media_type"
	ErrRangeNotSatisfiable         ErrCode = "range_not_satisfiable"
	ErrInsufficientStorage         ErrCode = "insufficient_storage"
	ErrLoopDetected                ErrCode = "loop_detected"
	ErrNotAcceptable               ErrCode = "not_acceptable"
	ErrTooEarly                    ErrCode = "too_early"
	ErrRequestHeaderFieldsTooLarge ErrCode = "request_header_fields_too_large"
)

func CheckCode(err error, code ErrCode) bool {
	if err == nil {
		return false
	}
	if serr, ok := err.(*Err); ok {
		return serr.code == code
	}
	return false
}
