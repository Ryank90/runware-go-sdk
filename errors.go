package runware

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
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

// APIError represents an error returned by the Runware API with full context
type APIError struct {
	Message     string
	ErrorID     string
	TaskUUID    string
	TaskType    string
	RawResponse string // Raw JSON response from API for debugging
	Timestamp   time.Time
}

// Error implements the error interface with comprehensive error details
func (e *APIError) Error() string {
	var parts []string

	// Build a detailed, developer-friendly error message
	if e.TaskType != "" {
		parts = append(parts, fmt.Sprintf("Task: %s", e.TaskType))
	}

	if e.ErrorID != "" {
		parts = append(parts, fmt.Sprintf("Code: %s", e.ErrorID))
	}

	if e.Message != "" {
		parts = append(parts, fmt.Sprintf("Message: %s", e.Message))
	} else {
		parts = append(parts, "Message: (empty error message from API)")
	}

	if e.TaskUUID != "" {
		parts = append(parts, fmt.Sprintf("TaskUUID: %s", e.TaskUUID))
	}

	// If we have a raw response and the message was empty, include it
	if e.Message == "" && e.RawResponse != "" {
		parts = append(parts, fmt.Sprintf("Raw Response: %s", e.RawResponse))
	}

	if len(parts) == 0 {
		return "API error: unknown error (no details provided by API)"
	}

	return fmt.Sprintf("Runware API Error - %s", strings.Join(parts, " | "))
}

// IsRetryable returns whether this error might be resolved by retrying
func (e *APIError) IsRetryable() bool {
	// Common retryable error codes/IDs
	retryableErrors := []string{
		"rateLimitExceeded",
		"serviceUnavailable",
		"timeout",
		"temporaryError",
	}

	for _, retryable := range retryableErrors {
		if strings.Contains(strings.ToLower(e.ErrorID), strings.ToLower(retryable)) {
			return true
		}
	}

	return false
}

// NewAPIError creates a new APIError from an ErrorResponse with full context
func NewAPIError(errResp *ErrorResponse) *APIError {
	// Handle different error response formats
	message := errResp.Error
	if message == "" && errResp.Message != "" {
		message = errResp.Message
	}

	errorID := errResp.ErrorID
	if errorID == "" && errResp.Code != "" {
		errorID = errResp.Code
	}

	// Capture raw response for debugging
	rawJSON, _ := json.Marshal(errResp)

	return &APIError{
		Message:     message,
		ErrorID:     errorID,
		TaskUUID:    errResp.TaskUUID,
		TaskType:    errResp.TaskType,
		RawResponse: string(rawJSON),
		Timestamp:   time.Now(),
	}
}

// IsAPIError checks if an error is an APIError
func IsAPIError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr)
}

// TimeoutError provides detailed context about what timed out
type TimeoutError struct {
	TaskType      string
	TaskUUID      string
	Duration      time.Duration
	ExpectedCount int
	ReceivedCount int
}

// Error implements the error interface
func (e *TimeoutError) Error() string {
	if e.ExpectedCount > 1 {
		return fmt.Sprintf(
			"timeout after %v waiting for %s (TaskUUID: %s) - received %d/%d results",
			e.Duration, e.TaskType, e.TaskUUID, e.ReceivedCount, e.ExpectedCount,
		)
	}
	return fmt.Sprintf(
		"timeout after %v waiting for %s (TaskUUID: %s) - no response received",
		e.Duration, e.TaskType, e.TaskUUID,
	)
}

// IsTimeout checks if an error is a timeout error
func IsTimeout(err error) bool {
	if errors.Is(err, ErrTimeout) {
		return true
	}
	var timeoutErr *TimeoutError
	return errors.As(err, &timeoutErr)
}
