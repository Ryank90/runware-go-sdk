package runware

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestErrorConstants(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrNotConnected", ErrNotConnected},
		{"ErrAlreadyConnected", ErrAlreadyConnected},
		{"ErrConnectionClosed", ErrConnectionClosed},
		{"ErrInvalidAPIKey", ErrInvalidAPIKey},
		{"ErrTimeout", ErrTimeout},
		{"ErrInvalidRequest", ErrInvalidRequest},
		{"ErrInvalidResponse", ErrInvalidResponse},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("%s is nil", tt.name)
			}
			if tt.err.Error() == "" {
				t.Errorf("%s has empty error message", tt.name)
			}
		})
	}
}

func TestAPIErrorError(t *testing.T) {
	tests := []struct {
		name    string
		apiErr  *APIError
		wantStr bool
	}{
		{
			name: "with error ID",
			apiErr: &APIError{
				Message:  "Test error",
				ErrorID:  "err-123",
				TaskUUID: "task-456",
			},
			wantStr: true,
		},
		{
			name: "without error ID",
			apiErr: &APIError{
				Message: "Test error",
			},
			wantStr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := tt.apiErr.Error()
			if tt.wantStr && errStr == "" {
				t.Error("Error() returned empty string")
			}
			if tt.wantStr && errStr != "" && len(errStr) < len(tt.apiErr.Message) {
				t.Errorf("Error() returned string shorter than message: %v", errStr)
			}
		})
	}
}

func TestNewAPIError(t *testing.T) {
	errResp := &ErrorResponse{
		Error:    "Test error message",
		ErrorID:  "err-789",
		TaskUUID: "task-101",
		TaskType: "imageInference",
	}

	apiErr := NewAPIError(errResp)

	if apiErr.Message != errResp.Error {
		t.Errorf("Message = %v, want %v", apiErr.Message, errResp.Error)
	}

	if apiErr.ErrorID != errResp.ErrorID {
		t.Errorf("ErrorID = %v, want %v", apiErr.ErrorID, errResp.ErrorID)
	}

	if apiErr.TaskUUID != errResp.TaskUUID {
		t.Errorf("TaskUUID = %v, want %v", apiErr.TaskUUID, errResp.TaskUUID)
	}

	if apiErr.TaskType != errResp.TaskType {
		t.Errorf("TaskType = %v, want %v", apiErr.TaskType, errResp.TaskType)
	}
}

func TestIsAPIError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "is APIError",
			err:  &APIError{Message: "test"},
			want: true,
		},
		{
			name: "is not APIError",
			err:  errors.New("regular error"),
			want: false,
		},
		{
			name: "is nil",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAPIError(tt.err); got != tt.want {
				t.Errorf("IsAPIError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIErrorWrapping(t *testing.T) {
	apiErr := &APIError{Message: "test error"}
	wrappedErr := errors.New("wrapped: " + apiErr.Error())

	if !errors.Is(wrappedErr, wrappedErr) {
		t.Error("Error wrapping failed")
	}
}

func TestAPIErrorEnhanced(t *testing.T) {
	tests := []struct {
		name     string
		errResp  *ErrorResponse
		contains []string
	}{
		{
			name: "full error details",
			errResp: &ErrorResponse{
				Error:    "unsupported dimensions",
				ErrorID:  "ERR001",
				TaskUUID: "task-123",
				TaskType: "imageInference",
			},
			contains: []string{
				"Task: imageInference",
				"Code: ERR001",
				"Message: unsupported dimensions",
				"TaskUUID: task-123",
			},
		},
		{
			name: "empty message shows raw response",
			errResp: &ErrorResponse{
				Error:    "",
				ErrorID:  "ERR002",
				TaskUUID: "task-456",
				TaskType: "videoInference",
			},
			contains: []string{
				"empty error message from API",
				"Raw Response:",
				"ERR002",
			},
		},
		{
			name: "alternative error format (code and message fields)",
			errResp: &ErrorResponse{
				Message:  "rate limit exceeded",
				Code:     "RATE_LIMIT",
				TaskUUID: "task-789",
				TaskType: "imageInference",
			},
			contains: []string{
				"Code: RATE_LIMIT",
				"Message: rate limit exceeded",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr := NewAPIError(tt.errResp)
			errStr := apiErr.Error()

			for _, expected := range tt.contains {
				if !strings.Contains(errStr, expected) {
					t.Errorf("Error string missing expected content.\nExpected to contain: %s\nGot: %s",
						expected, errStr)
				}
			}

			// Verify raw response is captured
			if apiErr.RawResponse == "" {
				t.Error("RawResponse should be populated")
			}

			// Verify timestamp is set
			if apiErr.Timestamp.IsZero() {
				t.Error("Timestamp should be set")
			}
		})
	}
}

func TestAPIErrorIsRetryable(t *testing.T) {
	tests := []struct {
		name      string
		errorID   string
		retryable bool
	}{
		{
			name:      "rate limit error is retryable",
			errorID:   "rateLimitExceeded",
			retryable: true,
		},
		{
			name:      "service unavailable is retryable",
			errorID:   "serviceUnavailable",
			retryable: true,
		},
		{
			name:      "timeout is retryable",
			errorID:   "timeout",
			retryable: true,
		},
		{
			name:      "validation error is not retryable",
			errorID:   "invalidParameter",
			retryable: false,
		},
		{
			name:      "auth error is not retryable",
			errorID:   "authenticationFailed",
			retryable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr := &APIError{
				ErrorID: tt.errorID,
			}

			if got := apiErr.IsRetryable(); got != tt.retryable {
				t.Errorf("IsRetryable() = %v, want %v", got, tt.retryable)
			}
		})
	}
}

func TestTimeoutError(t *testing.T) {
	tests := []struct {
		name     string
		err      *TimeoutError
		contains []string
	}{
		{
			name: "single result timeout",
			err: &TimeoutError{
				TaskType:      "imageInference",
				TaskUUID:      "task-123",
				Duration:      2 * time.Minute,
				ExpectedCount: 1,
				ReceivedCount: 0,
			},
			contains: []string{
				"2m0s",
				"imageInference",
				"task-123",
				"no response received",
			},
		},
		{
			name: "multiple results partial timeout",
			err: &TimeoutError{
				TaskType:      "imageInference",
				TaskUUID:      "task-456",
				Duration:      90 * time.Second,
				ExpectedCount: 4,
				ReceivedCount: 2,
			},
			contains: []string{
				"1m30s",
				"imageInference",
				"task-456",
				"received 2/4 results",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errStr := tt.err.Error()

			for _, expected := range tt.contains {
				if !strings.Contains(errStr, expected) {
					t.Errorf("Error string missing expected content.\nExpected to contain: %s\nGot: %s",
						expected, errStr)
				}
			}
		})
	}
}

func TestIsTimeout(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		isTimeout bool
	}{
		{
			name:      "TimeoutError is timeout",
			err:       &TimeoutError{TaskType: "test"},
			isTimeout: true,
		},
		{
			name:      "ErrTimeout is timeout",
			err:       ErrTimeout,
			isTimeout: true,
		},
		{
			name:      "APIError is not timeout",
			err:       &APIError{},
			isTimeout: false,
		},
		{
			name:      "nil is not timeout",
			err:       nil,
			isTimeout: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTimeout(tt.err); got != tt.isTimeout {
				t.Errorf("IsTimeout() = %v, want %v", got, tt.isTimeout)
			}
		})
	}
}

func TestDebugLogger(t *testing.T) {
	// Test default logger (no-op)
	defaultLog := &defaultLogger{}
	defaultLog.Printf("test message") // Should not panic

	// Test std logger
	stdLog := &stdLogger{}
	stdLog.Printf("test message %s", "arg") // Should not panic
}
