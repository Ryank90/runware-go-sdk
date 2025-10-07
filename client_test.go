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

	if req.TaskType != "upscaleGan" {
		t.Errorf("TaskType = %v, want upscaleGan", req.TaskType)
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
