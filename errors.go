package runware

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	models "github.com/Ryank90/runware-go-sdk/models"
)

// Common errors returned by the SDK.
// Use errors.Is() to check for these sentinel errors in your error handling.
var (
	// ErrNotConnected is returned when attempting to use a client that is not connected.
	// Call Connect() before making API requests.
	ErrNotConnected = errors.New("client is not connected")

	// ErrAlreadyConnected is returned when attempting to connect an already connected client.
	// This is typically not an error condition - just continue using the existing connection.
	ErrAlreadyConnected = errors.New("client is already connected")

	// ErrConnectionClosed is returned when the WebSocket connection is closed unexpectedly.
	// If auto-reconnect is enabled, the SDK will attempt to reconnect automatically.
	ErrConnectionClosed = errors.New("connection closed")

	// ErrInvalidAPIKey is returned when the API key is invalid or missing.
	// Ensure RUNWARE_API_KEY is set or provide it in the Config.
	ErrInvalidAPIKey = errors.New("invalid or missing API key")

	// ErrTimeout is returned when an operation times out.
	// Consider increasing RequestTimeout in the Config for longer-running operations.
	ErrTimeout = errors.New("operation timed out")

	// ErrInvalidRequest is returned when the request parameters are invalid.
	// Check that all required fields are populated.
	ErrInvalidRequest = errors.New("invalid request")

	// ErrInvalidResponse is returned when the API response cannot be parsed.
	// This may indicate an API version mismatch or network corruption.
	ErrInvalidResponse = errors.New("invalid response")
)

// APIError represents an error returned by the Runware API with full context.
//
// APIError provides detailed information about what went wrong, including:
//   - Message: Human-readable error description
//   - ErrorID: Machine-readable error code for programmatic handling
//   - TaskUUID: Unique identifier for the failed task (useful for support)
//   - TaskType: Type of operation that failed
//   - RawResponse: Full JSON response for debugging
//   - Timestamp: When the error occurred
//
// Use IsRetryable() to determine if the error is transient and worth retrying.
//
// Example:
//
//	if apiErr, ok := err.(*runware.APIError); ok {
//	    fmt.Printf("Error: %s (Code: %s)\n", apiErr.Message, apiErr.ErrorID)
//	    if apiErr.IsRetryable() {
//	        // Implement exponential backoff retry
//	    }
//	}
type APIError struct {
	// Message is the human-readable error description from the API
	Message string
	// ErrorID is the machine-readable error code (e.g., "rateLimitExceeded")
	ErrorID string
	// TaskUUID is the unique identifier for the failed task
	TaskUUID string
	// TaskType is the type of operation that failed (e.g., "imageInference")
	TaskType string
	// RawResponse is the full JSON response from the API for debugging
	RawResponse string
	// Timestamp is when the error occurred (client-side)
	Timestamp time.Time
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

// IsRetryable returns whether this error might be resolved by retrying the request.
//
// Returns true for transient errors like rate limiting, temporary service
// unavailability, or timeouts. Returns false for permanent errors like
// validation failures or authentication errors.
//
// Example retry logic:
//
//	for attempt := 0; attempt < maxRetries; attempt++ {
//	    resp, err := client.TextToImage(ctx, prompt, model, width, height)
//	    if err == nil {
//	        return resp, nil
//	    }
//	    if apiErr, ok := err.(*runware.APIError); ok && apiErr.IsRetryable() {
//	        time.Sleep(time.Duration(attempt+1) * time.Second)  // Exponential backoff
//	        continue
//	    }
//	    return nil, err  // Non-retryable error
//	}
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
func NewAPIError(errResp *models.ErrorResponse) *APIError {
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

// TimeoutError provides detailed context about timeout errors.
//
// This error is returned when an API request exceeds the configured timeout
// period. It includes information about what operation timed out and, for
// batch operations, how many results were received before timing out.
//
// Use IsTimeout() to check if an error is a timeout:
//
//	if runware.IsTimeout(err) {
//	    // Handle timeout - maybe increase timeout duration
//	}
type TimeoutError struct {
	// TaskType is the type of operation that timed out
	TaskType string
	// TaskUUID is the unique identifier for the task
	TaskUUID string
	// Duration is how long we waited before timing out
	Duration time.Duration
	// ExpectedCount is how many results we were expecting (for batch operations)
	ExpectedCount int
	// ReceivedCount is how many results we actually received before timeout
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
