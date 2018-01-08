package gremlin

import "errors"

const (
	StatusSuccess                  = 200
	StatusNoContent                = 204
	StatusPartialContent           = 206
	StatusUnauthorized             = 401
	StatusAuthenticate             = 407
	StatusMalformedRequest         = 498
	StatusInvalidRequestArguments  = 499
	StatusServerError              = 500
	StatusScriptEvaluationError    = 597
	StatusServerTimeout            = 598
	StatusServerSerializationError = 599
)

var (
	ErrStatusUnauthorized             = errors.New("Unauthorized")
	ErrStatusAuthenticate             = errors.New("Authenticate")
	ErrStatusMalformedRequest         = errors.New("Malformed Request")
	ErrStatusInvalidRequestArguments  = errors.New("Invalid Request Arguments")
	ErrStatusServerError              = errors.New("Server Error")
	ErrStatusScriptEvaluationError    = errors.New("Script Evaluation Error")
	ErrStatusServerTimeout            = errors.New("Server Timeout")
	ErrStatusServerSerializationError = errors.New("Server Serialization Error")
)

var ErrorMsg = map[int]error{
	StatusUnauthorized:             ErrStatusUnauthorized,
	StatusAuthenticate:             ErrStatusAuthenticate,
	StatusMalformedRequest:         ErrStatusMalformedRequest,
	StatusInvalidRequestArguments:  ErrStatusInvalidRequestArguments,
	StatusServerError:              ErrStatusServerError,
	StatusScriptEvaluationError:    ErrStatusScriptEvaluationError,
	StatusServerTimeout:            ErrStatusServerTimeout,
	StatusServerSerializationError: ErrStatusServerSerializationError,
}
