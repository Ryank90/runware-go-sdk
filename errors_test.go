package runware

import (
	"errors"
	"testing"
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
