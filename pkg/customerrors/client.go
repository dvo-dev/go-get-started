package customerrors

import (
	"fmt"
	"net/http"
)

type ClientErrorMessage struct {
	Error string `json:"error"`
}

type ClientError interface {
	// Error returns a human readable explanation of the error (logging)
	Error() string

	// StatusCode determines the http status code sent to the client
	StatusCode() int

	// ClientErrorMsg returns a ClientErrorMessage to be JSONified to the client
	ClientErrorMsg() ClientErrorMessage
}

type ClientErrorBadMethod struct {
	RequestMethod string
}

func (e ClientErrorBadMethod) Error() string {
	return fmt.Sprintf(
		`request method: '%s' not supported`,
		e.RequestMethod,
	)
}

func (e ClientErrorBadMethod) StatusCode() int {
	return http.StatusBadRequest
}

func (e ClientErrorBadMethod) ClientErrorMsg() ClientErrorMessage {
	return ClientErrorMessage{
		Error: e.Error(),
	}
}
