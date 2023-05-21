package helpers

import "fmt"

type Error interface {
	handleError(err error)
}

type ErrorHandler struct {
	message    string `json:"message"`
	err        error `json:"error"`
	errType   Error `json:"errorType"`
}

func WrapErrorf(err error, errType Error, format string, a ...interface{}) *ErrorHandler {
	return &ErrorHandler{
		message:    fmt.Sprintf(format, a...),
		err:        err,
		errType:    errType,
	}
}

func NewErrorf(errType Error, format string, a ...interface{}) *ErrorHandler {
	return WrapErrorf(nil, errType, format, a...)
}

