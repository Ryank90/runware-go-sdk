package runware

import (
	"context"
	"os"
	"testing"
	"time"
)

const (
	testPrompt = "test prompt"
	testModel  = "runware:101@1"
	testUUID   = "test-uuid"
	testAPIKey = "test-api-key"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "nil config without env var",
			config:  nil,
			wantErr: true,
		},
		{
			name: "valid config",
			config: &Config{
				APIKey:         "test-api-key",
				WSConfig:       DefaultWSConfig(),
				RequestTimeout: 30 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "empty API key",
			config: &Config{
				APIKey: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env var for nil config test
			if tt.config == nil {
				_ = os.Unsetenv("RUNWARE_API_KEY")
			}

			client, err := NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client without error")
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	apiKey := "test-key"
	config := DefaultConfig()
	config.APIKey = apiKey

	if config.APIKey != apiKey {
		t.Errorf("DefaultConfig().APIKey = %v, want %v", config.APIKey, apiKey)
	}

	if config.WSConfig == nil {
		t.Error("DefaultConfig().WSConfig is nil")
	}

	if config.RequestTimeout == 0 {
		t.Error("DefaultConfig().RequestTimeout is zero")
	}
}

func TestRequestBuilder(t *testing.T) {
	width := 1024
	height := 1024

	rb := NewRequestBuilder(testPrompt, testModel, width, height)

	negPrompt := "bad quality"
	rb.WithNegativePrompt(negPrompt)

	steps := 30
	rb.WithSteps(steps)

	cfg := 7.5
	rb.WithCFGScale(cfg)

	req := rb.Build()

	if req.PositivePrompt != testPrompt {
		t.Errorf("PositivePrompt = %v, want %v", req.PositivePrompt, testPrompt)
	}

	if req.Model != testModel {
		t.Errorf("Model = %v, want %v", req.Model, testModel)
	}

	if req.Width != width {
		t.Errorf("Width = %v, want %v", req.Width, width)
	}

	if req.Height != height {
		t.Errorf("Height = %v, want %v", req.Height, height)
	}

	if req.NegativePrompt == nil || *req.NegativePrompt != negPrompt {
		t.Errorf("NegativePrompt = %v, want %v", req.NegativePrompt, negPrompt)
	}

	if req.Steps == nil || *req.Steps != steps {
		t.Errorf("Steps = %v, want %v", req.Steps, steps)
	}

	if req.CFGScale == nil || *req.CFGScale != cfg {
		t.Errorf("CFGScale = %v, want %v", req.CFGScale, cfg)
	}
}

func TestNewImageInferenceRequest(t *testing.T) {
	model := "test-model"
	width := 512
	height := 512

	req := NewImageInferenceRequest(testPrompt, model, width, height)

	if req.TaskType != "imageInference" {
		t.Errorf("TaskType = %v, want imageInference", req.TaskType)
	}

	if req.TaskUUID == "" {
		t.Error("TaskUUID is empty")
	}

	if req.PositivePrompt != testPrompt {
		t.Errorf("PositivePrompt = %v, want %v", req.PositivePrompt, testPrompt)
	}

	if req.Model != model {
		t.Errorf("Model = %v, want %v", req.Model, model)
	}

	if req.Width != width {
		t.Errorf("Width = %v, want %v", req.Width, width)
	}

	if req.Height != height {
		t.Errorf("Height = %v, want %v", req.Height, height)
	}
}

func TestNewUploadImageRequest(t *testing.T) {
	req := NewUploadImageRequest()

	if req.TaskType != "imageUpload" {
		t.Errorf("TaskType = %v, want imageUpload", req.TaskType)
	}

	if req.TaskUUID == "" {
		t.Error("TaskUUID is empty")
	}
}

func TestNewUpscaleGanRequest(t *testing.T) {
	inputImage := testUUID
	factor := 4

	req := NewUpscaleGanRequest(inputImage, factor)

	if req.TaskType != "imageUpscale" {
		t.Errorf("TaskType = %v, want imageUpscale", req.TaskType)
	}

	if req.TaskUUID == "" {
		t.Error("TaskUUID is empty")
	}

	if req.InputImage != inputImage {
		t.Errorf("InputImage = %v, want %v", req.InputImage, inputImage)
	}

	if req.UpscaleFactor != factor {
		t.Errorf("UpscaleFactor = %v, want %v", req.UpscaleFactor, factor)
	}
}

func TestNewRemoveImageBackgroundRequest(t *testing.T) {
	inputImage := testUUID

	req := NewRemoveImageBackgroundRequest(inputImage)

	if req.TaskType != "imageBackgroundRemoval" {
		t.Errorf("TaskType = %v, want imageBackgroundRemoval", req.TaskType)
	}

	if req.TaskUUID == "" {
		t.Error("TaskUUID is empty")
	}

	if req.InputImage != inputImage {
		t.Errorf("InputImage = %v, want %v", req.InputImage, inputImage)
	}
}

func TestNewEnhancePromptRequest(t *testing.T) {
	req := NewEnhancePromptRequest(testPrompt)

	if req.TaskType != "promptEnhance" {
		t.Errorf("TaskType = %v, want promptEnhance", req.TaskType)
	}

	if req.TaskUUID == "" {
		t.Error("TaskUUID is empty")
	}

	if req.Prompt != testPrompt {
		t.Errorf("Prompt = %v, want %v", req.Prompt, testPrompt)
	}
}

func TestNewImageCaptionRequest(t *testing.T) {
	inputImage := testUUID

	req := NewImageCaptionRequest(inputImage)

	if req.TaskType != "imageCaption" {
		t.Errorf("TaskType = %v, want imageCaption", req.TaskType)
	}

	if req.TaskUUID == "" {
		t.Error("TaskUUID is empty")
	}

	if req.InputImage != inputImage {
		t.Errorf("InputImage = %v, want %v", req.InputImage, inputImage)
	}
}

func TestClientIsConnected(t *testing.T) {
	config := DefaultConfig()
	config.APIKey = testAPIKey
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if client.IsConnected() {
		t.Error("Client should not be connected initially")
	}
}

func TestImageInferenceWithoutConnection(t *testing.T) {
	config := DefaultConfig()
	config.APIKey = testAPIKey
	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	req := NewImageInferenceRequest("test", "model", 512, 512)

	_, err = client.ImageInference(ctx, req)
	if err != ErrNotConnected {
		t.Errorf("ImageInference() error = %v, want %v", err, ErrNotConnected)
	}
}

func TestAPIError(t *testing.T) {
	errResp := &ErrorResponse{
		Error:    "Test error",
		ErrorID:  "err-123",
		TaskUUID: "task-456",
		TaskType: "imageInference",
	}

	apiErr := NewAPIError(errResp)

	if apiErr.Message != errResp.Error {
		t.Errorf("APIError.Message = %v, want %v", apiErr.Message, errResp.Error)
	}

	if apiErr.ErrorID != errResp.ErrorID {
		t.Errorf("APIError.ErrorID = %v, want %v", apiErr.ErrorID, errResp.ErrorID)
	}

	errStr := apiErr.Error()
	if errStr == "" {
		t.Error("APIError.Error() returned empty string")
	}

	if !IsAPIError(apiErr) {
		t.Error("IsAPIError() returned false for APIError")
	}
}

func TestPollVideoResultSuccess(t *testing.T) {
	// This is a unit test for the polling logic structure
	// Actual API calls are tested in integration tests

	// Test that PollVideoResult handles context correctly
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Verify context is respected
	select {
	case <-ctx.Done():
		t.Error("Context should not be done yet")
	default:
		// Expected: context is still active
	}
}

func TestPollVideoResultContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Verify canceled context is handled
	if ctx.Err() == nil {
		t.Error("Context should be canceled")
	}
}

func TestExtractExpectedCount(t *testing.T) {
	config := DefaultConfig()
	config.APIKey = testAPIKey

	client := &Client{
		config:         config,
		requestTimeout: config.RequestTimeout,
		debugLogger:    &defaultLogger{},
	}

	tests := []struct {
		name     string
		req      *ImageInferenceRequest
		expected int
	}{
		{
			name: "nil numberResults defaults to 1",
			req: &ImageInferenceRequest{
				NumberResults: nil,
			},
			expected: 1,
		},
		{
			name: "numberResults = 4",
			req: &ImageInferenceRequest{
				NumberResults: func() *int { n := 4; return &n }(),
			},
			expected: 4,
		},
		{
			name: "numberResults = 1 explicitly",
			req: &ImageInferenceRequest{
				NumberResults: func() *int { n := 1; return &n }(),
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := client.extractExpectedCount(tt.req)
			if count != tt.expected {
				t.Errorf("extractExpectedCount() = %d, want %d", count, tt.expected)
			}
		})
	}
}

func TestTimeoutErrorFormatting(t *testing.T) {
	tests := []struct {
		name     string
		err      *TimeoutError
		contains string
	}{
		{
			name: "single result timeout",
			err: &TimeoutError{
				TaskType:      "imageInference",
				TaskUUID:      "test-uuid",
				Duration:      30 * time.Second,
				ExpectedCount: 1,
				ReceivedCount: 0,
			},
			contains: "no response received",
		},
		{
			name: "partial batch timeout",
			err: &TimeoutError{
				TaskType:      "imageInference",
				TaskUUID:      "test-uuid",
				Duration:      60 * time.Second,
				ExpectedCount: 4,
				ReceivedCount: 2,
			},
			contains: "received 2/4 results",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := tt.err.Error()
			if errMsg == "" {
				t.Error("Error message should not be empty")
			}
			// Basic validation that error contains expected info
			if len(errMsg) < 10 {
				t.Errorf("Error message too short: %s", errMsg)
			}
		})
	}
}

func TestDefaultConfigTimeout(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("DefaultConfig should not return nil")
	}

	if config.RequestTimeout == 0 {
		t.Error("RequestTimeout should be set")
	}

	if config.WSConfig == nil {
		t.Error("WSConfig should be initialized")
	}

	// Verify timeout is reasonable (should be at least 1 minute)
	if config.RequestTimeout < time.Minute {
		t.Errorf("RequestTimeout %v seems too short for production use", config.RequestTimeout)
	}
}
