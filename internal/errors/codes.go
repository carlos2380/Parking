package errors

import "net/http"

var (
	ErrInternalServerError    = &HttpError{http.StatusInternalServerError, "internal server error", nil}
	ErrNotFound               = &HttpError{http.StatusNotFound, "resource not found", nil}
	ErrBadRequest             = &HttpError{http.StatusBadRequest, "bad request", nil}
	ErrUnauthorized           = &HttpError{http.StatusUnauthorized, "unauthorized", nil}
	ErrForbidden              = &HttpError{http.StatusForbidden, "forbidden", nil}
	ErrIntConvert             = &HttpError{http.StatusInternalServerError, "unsupported int", nil}
	ErrRedisConnectionFailure = &HttpError{http.StatusInternalServerError, "unable to connect to Redis", nil}
	ErrRedisCarEntryFailed    = &HttpError{http.StatusInternalServerError, "failed to register car entry", nil}
	ErrMethodNotAllowed       = &HttpError{http.StatusMethodNotAllowed, "method not allowed", nil}
	ErrParkingClosed           = &HttpError{http.StatusInsufficientStorage, "Parking closed", nil}
)
