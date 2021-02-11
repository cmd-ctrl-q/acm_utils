package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type restErr struct {
	Message string        `json:"message"`
	Status  int           `json:"status"`
	Error   string        `json:"error"`
	Causes  []interface{} `json:"causes"`
}

// RestErr interface
type RestErr interface {
	GetMessage() string
	GetStatus() int
	GetError() string
	GetCauses() []interface{}
}

func (e restErr) GetMessage() string {
	return e.Message
}

func (e restErr) GetStatus() int {
	return e.Status
}

func (e restErr) GetError() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: [%v ]",
		e.Message, e.Status, e.Error, e.Causes)
}

func (e restErr) GetCauses() []interface{} {
	return e.Causes
}

// NewRestError is a custom error function
func NewRestError(message string, status int, err string, causes []interface{}) RestErr {
	return restErr{
		Message: message,
		Status:  status,
		Error:   err,
		Causes:  causes,
	}
}

// NewError returns a general message of the error.
// NewError is largely used to send a vague description back to an external caller.
func NewError(msg string) error {
	return errors.New(msg)
}

// NewRestErrorFromBytes attempts to create a RestErr
// If bytes object cannot be unmarshalled then return an
// invalid json error back to caller
func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr restErr
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}

// NewBadRequestError returns a status bad request
func NewBadRequestError(message string) RestErr {
	return restErr{
		Message: message,
		Status:  http.StatusBadRequest,
		Error:   "bad_request",
	}
}

// NewNotFoundError returns a status not found
func NewNotFoundError(message string) RestErr {
	return restErr{
		Message: message,
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}

// NewUnauthorizedError returns a rest error for unauthorized user
func NewUnauthorizedError(message string) RestErr {
	return restErr{
		Message: message,
		Status:  http.StatusUnauthorized,
		Error:   "unauthorized",
	}
}

// NewInternalServerError returns an internal server error
func NewInternalServerError(message string, err error) RestErr {
	result := restErr{
		Message: message,
		Status:  http.StatusInternalServerError,
		Error:   "internal_server_error",
	}
	if err != nil {
		result.Causes = append(result.Causes, err.Error())
	}
	return result
}
