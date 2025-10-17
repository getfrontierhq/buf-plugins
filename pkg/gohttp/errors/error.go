package errors

import (
	"errors"
	"net/http"
)

var _ error = &Error{}

type Error struct {
	Data            interface{} `json:"data"`
	Message         string      `json:"message"`
	InternalMessage string      `json:"-"`
}

// Error implements error.
func (e *Error) Error() string {
	if e.InternalMessage == "" {
		return e.Message
	}

	return e.Message + ": " + e.InternalMessage
}

func (e *Error) WithData(data interface{}) *Error {
	e.Data = data
	return e
}

func (e *Error) WithInternalMessage(internalMessage string) *Error {
	e.InternalMessage = internalMessage
	return e
}

func New(message string) *Error {
	return &Error{
		Message: message,
	}
}

var (
	ErrGeneralBadRequest                    = New("general error, Bad Request")
	ErrGeneralUnauthorized                  = New("general error, Unauthorized")
	ErrGeneralPaymentRequired               = New("general error, Payment Required")
	ErrGeneralForbidden                     = New("general error, Forbidden")
	ErrGeneralNotFound                      = New("general error, Not Found")
	ErrGeneralMethodNotAllowed              = New("general error, Method Not Allowed")
	ErrGeneralNotAcceptable                 = New("general error, Not Acceptable")
	ErrGeneralProxyAuthRequired             = New("general error, Proxy Authentication Required")
	ErrGeneralRequestTimeout                = New("general error, Request Timeout")
	ErrGeneralConflict                      = New("general error, Conflict")
	ErrGeneralGone                          = New("general error, Gone")
	ErrGeneralLengthRequired                = New("general error, Length Required")
	ErrGeneralPreconditionFailed            = New("general error, Precondition Failed")
	ErrGeneralRequestEntityTooLarge         = New("general error, Request Entity Too Large")
	ErrGeneralRequestURITooLong             = New("general error, Request URI Too Long")
	ErrGeneralUnsupportedMediaType          = New("general error, Unsupported Media Type")
	ErrGeneralRequestedRangeNotSatisfiable  = New("general error, Requested Range Not Satisfiable")
	ErrGeneralExpectationFailed             = New("general error, Expectation Failed")
	ErrGeneralTeapot                        = New("general error, I'm a teapot")
	ErrGeneralMisdirectedRequest            = New("general error, Misdirected Request")
	ErrGeneralUnprocessableEntity           = New("general error, Unprocessable Entity")
	ErrGeneralLocked                        = New("general error, Locked")
	ErrGeneralFailedDependency              = New("general error, Failed Dependency")
	ErrGeneralTooEarly                      = New("general error, Too Early")
	ErrGeneralUpgradeRequired               = New("general error, Upgrade Required")
	ErrGeneralPreconditionRequired          = New("general error, Precondition Required")
	ErrGeneralTooManyRequests               = New("general error, Too Many Requests")
	ErrGeneralRequestHeaderFieldsTooLarge   = New("general error, Request Header Fields Too Large")
	ErrGeneralUnavailableForLegalReasons    = New("general error, Unavailable For Legal Reasons")
	ErrGeneralInternalServerError           = New("general error, Internal Server Error")
	ErrGeneralNotImplemented                = New("general error, Not Implemented")
	ErrGeneralBadGateway                    = New("general error, Bad Gateway")
	ErrGeneralServiceUnavailable            = New("general error, Service Unavailable")
	ErrGeneralGatewayTimeout                = New("general error, Gateway Timeout")
	ErrGeneralHTTPVersionNotSupported       = New("general error, HTTP Version Not Supported")
	ErrGeneralVariantAlsoNegotiates         = New("general error, Variant Also Negotiates")
	ErrGeneralInsufficientStorage           = New("general error, Insufficient Storage")
	ErrGeneralLoopDetected                  = New("general error, Loop Detected")
	ErrGeneralNotExtended                   = New("general error, Not Extended")
	ErrGeneralNetworkAuthenticationRequired = New("general error, Network Authentication Required")
)

// Error map to associate status codes with error variables
var ErrorMap = map[int]error{
	http.StatusBadRequest:                    ErrGeneralBadRequest,
	http.StatusUnauthorized:                  ErrGeneralUnauthorized,
	http.StatusPaymentRequired:               ErrGeneralPaymentRequired,
	http.StatusForbidden:                     ErrGeneralForbidden,
	http.StatusNotFound:                      ErrGeneralNotFound,
	http.StatusMethodNotAllowed:              ErrGeneralMethodNotAllowed,
	http.StatusNotAcceptable:                 ErrGeneralNotAcceptable,
	http.StatusProxyAuthRequired:             ErrGeneralProxyAuthRequired,
	http.StatusRequestTimeout:                ErrGeneralRequestTimeout,
	http.StatusConflict:                      ErrGeneralConflict,
	http.StatusGone:                          ErrGeneralGone,
	http.StatusLengthRequired:                ErrGeneralLengthRequired,
	http.StatusPreconditionFailed:            ErrGeneralPreconditionFailed,
	http.StatusRequestEntityTooLarge:         ErrGeneralRequestEntityTooLarge,
	http.StatusRequestURITooLong:             ErrGeneralRequestURITooLong,
	http.StatusUnsupportedMediaType:          ErrGeneralUnsupportedMediaType,
	http.StatusRequestedRangeNotSatisfiable:  ErrGeneralRequestedRangeNotSatisfiable,
	http.StatusExpectationFailed:             ErrGeneralExpectationFailed,
	http.StatusTeapot:                        ErrGeneralTeapot,
	http.StatusMisdirectedRequest:            ErrGeneralMisdirectedRequest,
	http.StatusUnprocessableEntity:           ErrGeneralUnprocessableEntity,
	http.StatusLocked:                        ErrGeneralLocked,
	http.StatusFailedDependency:              ErrGeneralFailedDependency,
	http.StatusTooEarly:                      ErrGeneralTooEarly,
	http.StatusUpgradeRequired:               ErrGeneralUpgradeRequired,
	http.StatusPreconditionRequired:          ErrGeneralPreconditionRequired,
	http.StatusTooManyRequests:               ErrGeneralTooManyRequests,
	http.StatusRequestHeaderFieldsTooLarge:   ErrGeneralRequestHeaderFieldsTooLarge,
	http.StatusUnavailableForLegalReasons:    ErrGeneralUnavailableForLegalReasons,
	http.StatusInternalServerError:           ErrGeneralInternalServerError,
	http.StatusNotImplemented:                ErrGeneralNotImplemented,
	http.StatusBadGateway:                    ErrGeneralBadGateway,
	http.StatusServiceUnavailable:            ErrGeneralServiceUnavailable,
	http.StatusGatewayTimeout:                ErrGeneralGatewayTimeout,
	http.StatusHTTPVersionNotSupported:       ErrGeneralHTTPVersionNotSupported,
	http.StatusVariantAlsoNegotiates:         ErrGeneralVariantAlsoNegotiates,
	http.StatusInsufficientStorage:           ErrGeneralInsufficientStorage,
	http.StatusLoopDetected:                  ErrGeneralLoopDetected,
	http.StatusNotExtended:                   ErrGeneralNotExtended,
	http.StatusNetworkAuthenticationRequired: ErrGeneralNetworkAuthenticationRequired,
}

func ParseErr(err error) (int, error) {
	switch {
	case errors.Is(err, ErrGeneralBadRequest):
		return http.StatusBadRequest, err
	case errors.Is(err, ErrGeneralUnauthorized):
		return http.StatusUnauthorized, err
	case errors.Is(err, ErrGeneralPaymentRequired):
		return http.StatusPaymentRequired, err
	case errors.Is(err, ErrGeneralForbidden):
		return http.StatusForbidden, err
	case errors.Is(err, ErrGeneralNotFound):
		return http.StatusNotFound, err
	case errors.Is(err, ErrGeneralMethodNotAllowed):
		return http.StatusMethodNotAllowed, err
	case errors.Is(err, ErrGeneralNotAcceptable):
		return http.StatusNotAcceptable, err
	case errors.Is(err, ErrGeneralProxyAuthRequired):
		return http.StatusProxyAuthRequired, err
	case errors.Is(err, ErrGeneralRequestTimeout):
		return http.StatusRequestTimeout, err
	case errors.Is(err, ErrGeneralConflict):
		return http.StatusConflict, err
	case errors.Is(err, ErrGeneralGone):
		return http.StatusGone, err
	case errors.Is(err, ErrGeneralLengthRequired):
		return http.StatusLengthRequired, err
	case errors.Is(err, ErrGeneralPreconditionFailed):
		return http.StatusPreconditionFailed, err
	case errors.Is(err, ErrGeneralRequestEntityTooLarge):
		return http.StatusRequestEntityTooLarge, err
	case errors.Is(err, ErrGeneralRequestURITooLong):
		return http.StatusRequestURITooLong, err
	case errors.Is(err, ErrGeneralUnsupportedMediaType):
		return http.StatusUnsupportedMediaType, err
	case errors.Is(err, ErrGeneralRequestedRangeNotSatisfiable):
		return http.StatusRequestedRangeNotSatisfiable, err
	case errors.Is(err, ErrGeneralExpectationFailed):
		return http.StatusExpectationFailed, err
	case errors.Is(err, ErrGeneralTeapot):
		return http.StatusTeapot, err
	case errors.Is(err, ErrGeneralMisdirectedRequest):
		return http.StatusMisdirectedRequest, err
	case errors.Is(err, ErrGeneralUnprocessableEntity):
		return http.StatusUnprocessableEntity, err
	case errors.Is(err, ErrGeneralLocked):
		return http.StatusLocked, err
	case errors.Is(err, ErrGeneralFailedDependency):
		return http.StatusFailedDependency, err
	case errors.Is(err, ErrGeneralTooEarly):
		return http.StatusTooEarly, err
	case errors.Is(err, ErrGeneralUpgradeRequired):
		return http.StatusUpgradeRequired, err
	case errors.Is(err, ErrGeneralPreconditionRequired):
		return http.StatusPreconditionRequired, err
	case errors.Is(err, ErrGeneralTooManyRequests):
		return http.StatusTooManyRequests, err
	case errors.Is(err, ErrGeneralRequestHeaderFieldsTooLarge):
		return http.StatusRequestHeaderFieldsTooLarge, err
	case errors.Is(err, ErrGeneralUnavailableForLegalReasons):
		return http.StatusUnavailableForLegalReasons, err
	case errors.Is(err, ErrGeneralInternalServerError):
		return http.StatusInternalServerError, err
	case errors.Is(err, ErrGeneralNotImplemented):
		return http.StatusNotImplemented, err
	case errors.Is(err, ErrGeneralBadGateway):
		return http.StatusBadGateway, err
	case errors.Is(err, ErrGeneralServiceUnavailable):
		return http.StatusServiceUnavailable, err
	case errors.Is(err, ErrGeneralGatewayTimeout):
		return http.StatusGatewayTimeout, err
	case errors.Is(err, ErrGeneralHTTPVersionNotSupported):
		return http.StatusHTTPVersionNotSupported, err
	case errors.Is(err, ErrGeneralVariantAlsoNegotiates):
		return http.StatusVariantAlsoNegotiates, err
	case errors.Is(err, ErrGeneralInsufficientStorage):
		return http.StatusInsufficientStorage, err
	case errors.Is(err, ErrGeneralLoopDetected):
		return http.StatusLoopDetected, err
	case errors.Is(err, ErrGeneralNotExtended):
		return http.StatusNotExtended, err
	case errors.Is(err, ErrGeneralNetworkAuthenticationRequired):
		return http.StatusNetworkAuthenticationRequired, err
	default:
		return http.StatusInternalServerError, err
	}
}
