package runware

import (
	"errors"
	"fmt"
)

var (
	// ErrNotConnected is returned when attempting to use a client that is not connected
	ErrNotConnected = errors.New("client is not connected")

	// ErrAlreadyConnected is returned when attempting to connect an already connected client
	ErrAlreadyConnected = errors.New("client is already connected")

	// ErrConnectionClosed is returned when the connection is closed unexpectedly
	ErrConnectionClosed = errors.New("connection closed")

	// ErrInvalidAPIKey is returned when the API key is invalid or missing
	ErrInvalidAPIKey = errors.New("invalid or missing API key")

	// ErrTimeout is returned when an operation times out
	ErrTimeout = errors.New("operation timed out")

	// ErrInvalidRequest is returned when the request is invalid
	ErrInvalidRequest = errors.New("invalid request")

	// ErrInvalidResponse is returned when the response is invalid
	ErrInvalidResponse = errors.New("invalid response")
)

// APIError represents an error returned by the Runware API
type APIError struct {
	Message  string
	ErrorID  string
	TaskUUID string
	TaskType string
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.ErrorID != "" {
		return fmt.Sprintf("API error [%s]: %s (task: %s)", e.ErrorID, e.Message, e.TaskUUID)
	}
	return fmt.Sprintf("API error: %s", e.Message)
}

// NewAPIError creates a new APIError from an ErrorResponse
func NewAPIError(errResp *ErrorResponse) *APIError {
	return &APIError{
		Message:  errResp.Error,
		ErrorID:  errResp.ErrorID,
		TaskUUID: errResp.TaskUUID,
		TaskType: errResp.TaskType,
	}
}

// IsAPIError checks if an error is an APIError
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}
