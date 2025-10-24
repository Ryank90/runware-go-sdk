// Package runware provides a Go SDK for the Runware AI platform.
//
// The SDK enables generation, transformation, and enhancement of images, videos, and audio
// using state-of-the-art AI models through a simple, type-safe interface.
//
// Basic usage:
//
//	client, err := runware.NewClient(nil) // Uses RUNWARE_API_KEY env var
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Disconnect()
//
//	if err := client.Connect(context.Background()); err != nil {
//	    log.Fatal(err)
//	}
//
//	resp, err := client.TextToImage(ctx, "a serene mountain landscape", "runware:101@1", 1024, 1024)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Image URL: %s\n", *resp.ImageURL)
package runware

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	wsinternal "github.com/Ryank90/runware-go-sdk/internal/ws"
	models "github.com/Ryank90/runware-go-sdk/models"
)

// DebugLogger is an interface for custom debug logging implementations.
// It allows users to integrate SDK logging with their existing logging infrastructure.
type DebugLogger interface {
	// Printf logs a formatted debug message.
	Printf(format string, v ...interface{})
}

// defaultLogger is a no-op logger
type defaultLogger struct{}

func (d *defaultLogger) Printf(format string, v ...interface{}) {}

// stdLogger wraps standard log package
type stdLogger struct{}

func (s *stdLogger) Printf(format string, v ...interface{}) {
	log.Printf("[Runware SDK] "+format, v...)
}

// Client is the main Runware SDK client that manages WebSocket connections
// and provides methods for image, video, and audio generation.
//
// Client instances are safe for concurrent use by multiple goroutines.
// A single Client can handle multiple simultaneous requests efficiently
// through multiplexed WebSocket communication.
type Client struct {
	ws             *wsinternal.Client
	apiKey         string
	config         *Config
	requestTimeout time.Duration
	debugLogger    DebugLogger
}

// Config contains client configuration options.
//
// Use DefaultConfig() to get a configuration with sensible defaults,
// then customize as needed.
type Config struct {
	// APIKey is the Runware API key.
	// If empty when calling NewClient, it will be read from the RUNWARE_API_KEY environment variable.
	APIKey string

	// WSConfig contains WebSocket-specific configuration (connection timeouts, reconnection settings, etc.).
	// If nil, DefaultWSConfig() will be used.
	WSConfig *wsinternal.WSConfig

	// RequestTimeout is the default timeout for API requests.
	// Individual requests may override this timeout using context.WithTimeout.
	// Default: 120 seconds (suitable for video/image generation).
	RequestTimeout time.Duration

	// EnableDebugLogging enables detailed debug logs useful for troubleshooting connection issues.
	// Can also be enabled by setting RUNWARE_DEBUG=1 environment variable.
	EnableDebugLogging bool

	// DebugLogger is a custom logger for debug output.
	// If nil and EnableDebugLogging is true, logs will be written to standard log output.
	DebugLogger DebugLogger
}

// DefaultConfig returns a client configuration with sensible defaults.
//
// The returned config uses:
//   - Empty APIKey (must be set explicitly or via RUNWARE_API_KEY env var)
//   - 120-second request timeout
//   - Default WebSocket configuration with auto-reconnect enabled
//   - Debug logging enabled if RUNWARE_DEBUG=1 environment variable is set
func DefaultConfig() *Config {
	// Check environment variable for debug mode
	debugEnabled := os.Getenv("RUNWARE_DEBUG") == "true" || os.Getenv("RUNWARE_DEBUG") == "1"

	return &Config{
		APIKey:             "",
		WSConfig:           wsinternal.DefaultWSConfig(),
		RequestTimeout:     120 * time.Second, // Generous timeout for image/video generation
		EnableDebugLogging: debugEnabled,
	}
}

// NewClient creates a new Runware client.
//
// If config is nil, the client will attempt to read the API key from the
// RUNWARE_API_KEY environment variable and use default configuration.
//
// Returns ErrInvalidAPIKey if no API key is provided or found.
//
// Example:
//
//	// Use environment variable
//	client, err := runware.NewClient(nil)
//
//	// Use explicit configuration
//	config := &runware.Config{
//	    APIKey: "your-api-key",
//	    RequestTimeout: 60 * time.Second,
//	}
//	client, err := runware.NewClient(config)
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		apiKey := os.Getenv("RUNWARE_API_KEY")
		if apiKey == "" {
			return nil, ErrInvalidAPIKey
		}
		config = DefaultConfig()
		config.APIKey = apiKey
	}

	if config.APIKey == "" {
		return nil, ErrInvalidAPIKey
	}

	// Initialize debug logger
	var debugLogger DebugLogger
	if config.EnableDebugLogging {
		if config.DebugLogger != nil {
			debugLogger = config.DebugLogger
		} else {
			debugLogger = &stdLogger{}
		}
	} else {
		debugLogger = &defaultLogger{}
	}

	client := &Client{
		apiKey:         config.APIKey,
		config:         config,
		requestTimeout: config.RequestTimeout,
		debugLogger:    debugLogger,
		ws:             wsinternal.NewClient(config.APIKey, config.WSConfig, debugLogger),
	}

	return client, nil
}

// Connect establishes a WebSocket connection to the Runware API.
//
// This method must be called before making any API requests. The connection
// uses WebSockets for efficient, bidirectional communication and supports
// automatic reconnection if configured.
//
// Returns an error if the connection cannot be established within the configured timeout.
//
// Example:
//
//	ctx := context.Background()
//	if err := client.Connect(ctx); err != nil {
//	    log.Fatal("Failed to connect:", err)
//	}
//	defer client.Disconnect()
func (c *Client) Connect(ctx context.Context) error { return c.ws.Connect(ctx) }

// Disconnect closes the WebSocket connection to the Runware API.
//
// This method should be called when the client is no longer needed to free
// up resources. It's safe to call Disconnect multiple times.
//
// Any in-flight requests will receive an error after disconnection.
func (c *Client) Disconnect() error { return c.ws.Disconnect() }

// IsConnected returns whether the client currently has an active connection
// to the Runware API.
//
// This method is safe to call from multiple goroutines.
func (c *Client) IsConnected() bool { return c.ws.IsConnected() }

// ImageInference performs AI-powered image generation with full control over generation parameters.
//
// This is the low-level method that accepts a complete ImageInferenceRequest with all optional
// parameters. For common use cases, consider using the convenience methods:
//   - TextToImage() for simple text-to-image generation
//   - ImageToImage() for image-to-image transformation
//   - Inpaint() for inpainting with masks
//   - Outpaint() for extending images beyond their borders
//
// The request must specify at least:
//   - PositivePrompt: Text description of the desired image
//   - Model: AI model to use (e.g., "runware:101@1")
//   - Width and Height: Output dimensions in pixels
//
// Returns ErrNotConnected if the client is not connected.
// Returns ErrInvalidRequest if the request is nil or invalid.
//
// Example:
//
//	req := models.NewImageInferenceRequest("mountain landscape", "runware:101@1", 1024, 1024)
//	steps := 30
//	req.Steps = &steps
//	resp, err := client.ImageInference(ctx, req)
func (c *Client) ImageInference(ctx context.Context, req *models.ImageInferenceRequest) (*models.ImageInferenceResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	// Set default task type if not set
	if req.TaskType == "" {
		req.TaskType = models.TaskTypeImageInference
	}

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*models.ImageInferenceResponse), nil
}

// ImageInferenceBatch performs multiple image inference requests in parallel with bounded concurrency.
//
// This method efficiently processes multiple image generation requests concurrently,
// automatically managing goroutines to prevent resource exhaustion. Results are returned
// in the same order as the input requests.
//
// The method uses a semaphore to limit concurrent requests, making it safe to process
// large batches without overwhelming system resources.
//
// Returns a slice of responses corresponding to each request. If any request fails,
// the error is returned and processing stops.
//
// Example:
//
//	requests := []*models.ImageInferenceRequest{
//	    models.NewImageInferenceRequest("sunset", "runware:101@1", 1024, 1024),
//	    models.NewImageInferenceRequest("ocean", "runware:101@1", 1024, 1024),
//	}
//	responses, err := client.ImageInferenceBatch(ctx, requests)
func (c *Client) ImageInferenceBatch(ctx context.Context, requests []*models.ImageInferenceRequest) ([]*models.ImageInferenceResponse, error) {
	return processBatch(ctx, requests, c.ImageInference)
}

// UploadImage uploads an image to Runware
func (c *Client) UploadImage(ctx context.Context, req *models.UploadImageRequest) (*models.UploadImageResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.TaskType == "" {
		req.TaskType = models.TaskTypeImageUpload
	}

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*models.UploadImageResponse), nil
}

// UploadImageFromFile uploads an image from a file path
func (c *Client) UploadImageFromFile(ctx context.Context, filePath string) (*models.UploadImageResponse, error) {
	data, err := os.ReadFile(filePath) // #nosec G304 - file path is provided by user for upload
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	base64Data := base64.StdEncoding.EncodeToString(data)
	req := models.NewUploadImageRequest()
	req.ImageBase64 = &base64Data

	return c.UploadImage(ctx, req)
}

// UploadImageFromURL uploads an image from a URL
func (c *Client) UploadImageFromURL(ctx context.Context, url string) (*models.UploadImageResponse, error) {
	req := models.NewUploadImageRequest()
	req.ImageURL = &url

	return c.UploadImage(ctx, req)
}

// UpscaleImage upscales an image using GAN
func (c *Client) UpscaleImage(ctx context.Context, req *models.UpscaleGanRequest) (*models.UpscaleGanResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.TaskType == "" {
		req.TaskType = models.TaskTypeUpscaleGan
	}

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*models.UpscaleGanResponse), nil
}

// RemoveBackground removes the background from an image
func (c *Client) RemoveBackground(ctx context.Context, req *models.RemoveImageBackgroundRequest) (*models.RemoveImageBackgroundResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.TaskType == "" {
		req.TaskType = models.TaskTypeImageBackgroundRemoval
	}

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*models.RemoveImageBackgroundResponse), nil
}

// EnhancePrompt enhances a text prompt
func (c *Client) EnhancePrompt(ctx context.Context, req *models.EnhancePromptRequest) (*models.EnhancePromptResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.TaskType == "" {
		req.TaskType = models.TaskTypePromptEnhance
	}

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*models.EnhancePromptResponse), nil
}

// CaptionImage generates a caption for an image
func (c *Client) CaptionImage(ctx context.Context, req *models.ImageCaptionRequest) (*models.ImageCaptionResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.TaskType == "" {
		req.TaskType = models.TaskTypeImageCaption
	}

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*models.ImageCaptionResponse), nil
}

// VideoInference performs video inference (async only - returns acknowledgment)
// For video generation, this returns quickly with just the taskUUID acknowledgment.
// Use PollVideoResult() or GetResponse() to retrieve the actual video result.
func (c *Client) VideoInference(ctx context.Context, req *models.VideoInferenceRequest) (*models.VideoInferenceResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.TaskType == "" {
		req.TaskType = models.TaskTypeVideoInference
	}

	// Video is async-only, so this will return quickly with acknowledgment
	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*models.VideoInferenceResponse), nil
}

// VideoInferenceBatch performs multiple video inference requests in parallel
func (c *Client) VideoInferenceBatch(ctx context.Context, requests []*models.VideoInferenceRequest) ([]*models.VideoInferenceResponse, error) {
	return processBatch(ctx, requests, c.VideoInference)
}

// sendRequest is a generic method to send a request and wait for response
func (c *Client) sendRequest(ctx context.Context, req interface{}) (interface{}, error) {
	if !c.IsConnected() {
		return nil, ErrNotConnected
	}

	expectedCount := c.extractExpectedCount(req)
	respChan := make(chan interface{}, expectedCount)
	errChan := make(chan error, 1)

	// Define cleanup to remove handler from websocket after final response
	var taskUUID, taskType string
	if ti, ok := req.(models.TaskIdentifiable); ok {
		taskUUID = ti.GetTaskUUID()
		taskType = ti.GetTaskType()
	}
	onDone := func() {
		if taskUUID != "" {
			c.ws.RemoveHandler(taskUUID)
		}
	}
	handler := c.createResponseHandler(expectedCount, respChan, errChan, onDone)

	// Extract task info for error reporting (already set above)

	c.debugLogger.Printf("Submitting request: %s (TaskUUID: %s, expecting %d results)",
		taskType, taskUUID, expectedCount)

	// Send the request
	if err := c.ws.Send(ctx, req, handler); err != nil {
		return nil, err
	}

	return c.waitForResponse(ctx, taskType, taskUUID, expectedCount, respChan, errChan)
}

// extractExpectedCount extracts the numberResults from a request
func (c *Client) extractExpectedCount(req interface{}) int {
	expectedCount := 1

	if rcp, ok := req.(models.ResultCountProvider); ok {
		if nr := rcp.GetNumberResults(); nr != nil && *nr > 0 {
			expectedCount = *nr
		}
		return expectedCount
	}

	return expectedCount
}

// createResponseHandler creates a handler function for collecting responses
func (c *Client) createResponseHandler(
	expectedCount int,
	respChan chan interface{},
	errChan chan error,
	onDone func(),
) func(interface{}, error) {
	var mu sync.Mutex
	receivedCount := 0

	return func(data interface{}, err error) {
		mu.Lock()
		defer mu.Unlock()

		if err != nil {
			select {
			case errChan <- err:
			default:
			}
			return
		}

		receivedCount++
		respChan <- data

		if receivedCount >= expectedCount {
			close(respChan)
			if onDone != nil {
				onDone()
			}
		}
	}
}

// waitForResponse waits for the expected number of responses or timeout
func (c *Client) waitForResponse(
	ctx context.Context,
	taskType, taskUUID string,
	expectedCount int,
	respChan chan interface{},
	errChan chan error,
) (interface{}, error) {
	timeout := c.requestTimeout
	if deadline, ok := ctx.Deadline(); ok {
		timeout = time.Until(deadline)
	}

	startTime := time.Now()
	timeoutTimer := time.After(timeout)

	c.debugLogger.Printf("Waiting for response: %s (TaskUUID: %s, timeout: %v)",
		taskType, taskUUID, timeout)

	// For single result, return immediately
	if expectedCount == 1 {
		result, err := c.waitForSingleResponse(ctx, respChan, errChan, timeoutTimer)
		if err != nil && IsTimeout(err) {
			// Return enhanced timeout error
			return nil, &TimeoutError{
				TaskType:      taskType,
				TaskUUID:      taskUUID,
				Duration:      time.Since(startTime),
				ExpectedCount: expectedCount,
				ReceivedCount: 0,
			}
		}
		return result, err
	}

	// For multiple results, wait for all
	return c.waitForMultipleResponses(
		ctx, taskType, taskUUID, expectedCount, startTime, respChan, errChan, timeoutTimer,
	)
}

// waitForSingleResponse waits for a single response
func (c *Client) waitForSingleResponse(
	ctx context.Context,
	respChan chan interface{},
	errChan chan error,
	timeoutTimer <-chan time.Time,
) (interface{}, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errChan:
		return nil, err
	case resp := <-respChan:
		return resp, nil
	case <-timeoutTimer:
		return nil, ErrTimeout
	}
}

// waitForMultipleResponses waits for multiple responses
func (c *Client) waitForMultipleResponses(
	ctx context.Context,
	taskType, taskUUID string,
	expectedCount int,
	startTime time.Time,
	respChan chan interface{},
	errChan chan error,
	timeoutTimer <-chan time.Time,
) (interface{}, error) {
	var firstResp interface{}
	receivedCount := 0

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-errChan:
			return nil, err
		case resp, ok := <-respChan:
			if firstResp == nil {
				firstResp = resp
			}
			receivedCount++
			c.debugLogger.Printf("Received %d/%d results for %s (TaskUUID: %s)",
				receivedCount, expectedCount, taskType, taskUUID)
			if !ok {
				return firstResp, nil
			}
		case <-timeoutTimer:
			// Return enhanced timeout error with partial results info
			return nil, &TimeoutError{
				TaskType:      taskType,
				TaskUUID:      taskUUID,
				Duration:      time.Since(startTime),
				ExpectedCount: expectedCount,
				ReceivedCount: receivedCount,
			}
		}
	}
}

// Helper methods for common operations

// TextToImage generates an image from a text prompt
func (c *Client) TextToImage(ctx context.Context, prompt, model string, width, height int) (*models.ImageInferenceResponse, error) {
	req := models.NewImageInferenceRequest(prompt, model, width, height)
	return c.ImageInference(ctx, req)
}

// TextToVideo generates a video from a text prompt
func (c *Client) TextToVideo(ctx context.Context, prompt, model string, duration int) (*models.VideoInferenceResponse, error) {
	req := models.NewVideoInferenceRequest(prompt, model)
	req.Duration = &duration
	return c.VideoInference(ctx, req)
}

// ImageToVideo generates a video from an image and prompt
func (c *Client) ImageToVideo(ctx context.Context, prompt, model, seedImage string, duration int) (*models.VideoInferenceResponse, error) {
	req := models.NewVideoInferenceRequest(prompt, model)
	req.Duration = &duration
	req.FrameImages = []models.FrameImage{{InputImage: seedImage, Frame: models.FramePositionFirst}}
	return c.VideoInference(ctx, req)
}

// GetResponse polls for the result of an async task
// Returns either *VideoInferenceResponse or *AudioInferenceResponse depending on the task
func (c *Client) GetResponse(ctx context.Context, taskUUID string) (interface{}, error) {
	req := models.NewGetResponseRequest(taskUUID)

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// PollVideoResult polls for video generation result with automatic retries
// This function will keep polling until it gets a definitive response (success/error)
// or reaches maxAttempts. Each individual poll has its own timeout, but the overall
// polling process can run as long as needed.
func (c *Client) PollVideoResult(
	ctx context.Context,
	taskUUID string,
	maxAttempts int,
	pollInterval time.Duration,
) (*models.VideoInferenceResponse, error) {
	c.debugLogger.Printf("Starting video polling for TaskUUID: %s (max attempts: %d, interval: %v)",
		taskUUID, maxAttempts, pollInterval)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Check if parent context was canceled
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Create a fresh context for each poll with a reasonable timeout
		// This allows the overall polling to continue beyond the client's default timeout
		pollCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		resp, err := c.GetResponse(pollCtx, taskUUID)
		cancel()

		if err != nil {
			c.debugLogger.Printf("Poll attempt %d/%d: error - %v", attempt, maxAttempts, err)
			// If it's a timeout on this specific poll, continue trying
			if IsTimeout(err) || err == context.DeadlineExceeded {
				c.debugLogger.Printf("Poll request timed out, will retry...")
				time.Sleep(pollInterval)
				continue
			}
			// For other errors, return immediately
			return nil, err
		}

		// Type assert to VideoInferenceResponse
		videoResp, ok := resp.(*models.VideoInferenceResponse)
		if !ok {
			return nil, fmt.Errorf("unexpected response type from getResponse")
		}

		c.debugLogger.Printf("Poll attempt %d/%d: status = %s", attempt, maxAttempts, videoResp.Status)

		switch videoResp.Status {
		case models.TaskStatusSuccess:
			c.debugLogger.Printf("Video generation completed successfully!")
			return videoResp, nil
		case models.TaskStatusError:
			return nil, fmt.Errorf("video generation failed - check API response for details")
		case models.TaskStatusProcessing, "":
			// Continue polling - status is still processing or not set yet
			c.debugLogger.Printf("Video still processing, waiting %v before next poll...", pollInterval)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(pollInterval):
				continue
			}
		default:
			// Unknown status, log it but continue polling
			c.debugLogger.Printf("Unknown status: %s, continuing to poll...", videoResp.Status)
			time.Sleep(pollInterval)
		}
	}

	return nil, fmt.Errorf("polling exhausted: reached max attempts (%d) without definitive response", maxAttempts)
}

// ImageToImage transforms an image based on a prompt
// seedImageURLOrUUID can be either an image URL or UUID
func (c *Client) ImageToImage(
	ctx context.Context,
	prompt, model, seedImageURLOrUUID string,
	width, height int,
	strength float64,
) (*models.ImageInferenceResponse, error) {
	req := models.NewImageInferenceRequest(prompt, model, width, height)
	req.SeedImage = &seedImageURLOrUUID
	req.Strength = &strength
	return c.ImageInference(ctx, req)
}

// Inpaint performs inpainting on an image
func (c *Client) Inpaint(
	ctx context.Context,
	prompt, model, seedImage, maskImage string,
	width, height int,
	strength float64,
) (*models.ImageInferenceResponse, error) {
	req := models.NewImageInferenceRequest(prompt, model, width, height)
	req.SeedImage = &seedImage
	req.MaskImage = &maskImage
	req.Strength = &strength
	return c.ImageInference(ctx, req)
}

// Outpaint performs outpainting on an image
func (c *Client) Outpaint(
	ctx context.Context,
	prompt, model, seedImage string,
	width, height int,
	outpaint *models.Outpaint,
) (*models.ImageInferenceResponse, error) {
	req := models.NewImageInferenceRequest(prompt, model, width, height)
	req.SeedImage = &seedImage
	req.Outpaint = outpaint
	return c.ImageInference(ctx, req)
}

// RequestBuilder provides a fluent interface for building requests
type RequestBuilder struct{ req *models.ImageInferenceRequest }

// NewRequestBuilder creates a new request builder
func NewRequestBuilder(prompt, model string, width, height int) *RequestBuilder {
	return &RequestBuilder{req: models.NewImageInferenceRequest(prompt, model, width, height)}
}

// WithNegativePrompt sets the negative prompt
func (rb *RequestBuilder) WithNegativePrompt(prompt string) *RequestBuilder {
	rb.req.NegativePrompt = &prompt
	return rb
}

// WithSeedImage sets the seed image
func (rb *RequestBuilder) WithSeedImage(imageUUID string) *RequestBuilder {
	rb.req.SeedImage = &imageUUID
	return rb
}

// WithMaskImage sets the mask image
func (rb *RequestBuilder) WithMaskImage(imageUUID string) *RequestBuilder {
	rb.req.MaskImage = &imageUUID
	return rb
}

// WithStrength sets the strength
func (rb *RequestBuilder) WithStrength(strength float64) *RequestBuilder {
	rb.req.Strength = &strength
	return rb
}

// WithSteps sets the number of steps
func (rb *RequestBuilder) WithSteps(steps int) *RequestBuilder {
	rb.req.Steps = &steps
	return rb
}

// WithCFGScale sets the CFG scale
func (rb *RequestBuilder) WithCFGScale(scale float64) *RequestBuilder {
	rb.req.CFGScale = &scale
	return rb
}

// WithSeed sets the seed
func (rb *RequestBuilder) WithSeed(seed int64) *RequestBuilder {
	rb.req.Seed = &seed
	return rb
}

// WithScheduler sets the scheduler
func (rb *RequestBuilder) WithScheduler(scheduler models.Scheduler) *RequestBuilder {
	rb.req.Scheduler = &scheduler
	return rb
}

// WithNumberResults sets the number of results
func (rb *RequestBuilder) WithNumberResults(num int) *RequestBuilder {
	rb.req.NumberResults = &num
	return rb
}

// WithOutputType sets the output type
func (rb *RequestBuilder) WithOutputType(outputType models.OutputType) *RequestBuilder {
	rb.req.OutputType = &outputType
	return rb
}

// WithOutputFormat sets the output format
func (rb *RequestBuilder) WithOutputFormat(format models.OutputFormat) *RequestBuilder {
	rb.req.OutputFormat = &format
	return rb
}

// WithLoRA adds a LoRA
func (rb *RequestBuilder) WithLoRA(model string, weight float64) *RequestBuilder {
	rb.req.LoRA = append(rb.req.LoRA, models.LoRA{
		Model:  model,
		Weight: &weight,
	})
	return rb
}

// WithControlNet adds a ControlNet
func (rb *RequestBuilder) WithControlNet(model, guideImage string, weight float64) *RequestBuilder {
	rb.req.ControlNet = append(rb.req.ControlNet, models.ControlNet{
		Model:      model,
		GuideImage: guideImage,
		Weight:     &weight,
	})
	return rb
}

// WithEmbedding adds an embedding
func (rb *RequestBuilder) WithEmbedding(model string, weight float64) *RequestBuilder {
	rb.req.Embeddings = append(rb.req.Embeddings, models.Embedding{
		Model:  model,
		Weight: &weight,
	})
	return rb
}

// WithIPAdapter adds an IP-Adapter
func (rb *RequestBuilder) WithIPAdapter(model, guideImage string, weight float64) *RequestBuilder {
	rb.req.IPAdapters = append(rb.req.IPAdapters, models.IPAdapter{
		Model:      model,
		GuideImage: guideImage,
		Weight:     &weight,
	})
	return rb
}

// WithOutpaint sets the outpaint parameters
func (rb *RequestBuilder) WithOutpaint(outpaint *models.Outpaint) *RequestBuilder {
	rb.req.Outpaint = outpaint
	return rb
}

// WithRefiner sets the refiner
func (rb *RequestBuilder) WithRefiner(model string, startStep int) *RequestBuilder {
	rb.req.Refiner = &models.Refiner{
		Model:     model,
		StartStep: &startStep,
	}
	return rb
}

// WithSafety enables safety checks
func (rb *RequestBuilder) WithSafety(mode models.SafetyMode) *RequestBuilder {
	checkContent := true
	rb.req.Safety = &models.Safety{
		CheckContent: checkContent,
		Mode:         mode,
	}
	return rb
}

// WithIncludeCost includes cost in the response
func (rb *RequestBuilder) WithIncludeCost(include bool) *RequestBuilder {
	rb.req.IncludeCost = &include
	return rb
}

// Build returns the built request
func (rb *RequestBuilder) Build() *models.ImageInferenceRequest {
	return rb.req
}

// VideoRequestBuilder provides a fluent interface for building video requests
type VideoRequestBuilder struct{ req *models.VideoInferenceRequest }

// NewVideoRequestBuilder creates a new video request builder
func NewVideoRequestBuilder(prompt, model string) *VideoRequestBuilder {
	return &VideoRequestBuilder{req: models.NewVideoInferenceRequest(prompt, model)}
}

// WithNegativePrompt sets the negative prompt
func (vb *VideoRequestBuilder) WithNegativePrompt(prompt string) *VideoRequestBuilder {
	vb.req.NegativePrompt = &prompt
	return vb
}

// WithDuration sets the video duration in seconds (1-10)
func (vb *VideoRequestBuilder) WithDuration(duration int) *VideoRequestBuilder {
	vb.req.Duration = &duration
	return vb
}

// WithResolution sets the video resolution
func (vb *VideoRequestBuilder) WithResolution(width, height int) *VideoRequestBuilder {
	vb.req.Width = &width
	vb.req.Height = &height
	return vb
}

// WithFPS sets the frames per second (15-60)
func (vb *VideoRequestBuilder) WithFPS(fps int) *VideoRequestBuilder {
	vb.req.FPS = &fps
	return vb
}

// WithFrameImage adds a frame image constraint (first or last frame)
func (vb *VideoRequestBuilder) WithFrameImage(imageUUID string, position models.FramePosition) *VideoRequestBuilder {
	vb.req.FrameImages = append(vb.req.FrameImages, models.FrameImage{
		InputImage: imageUUID,
		Frame:      position,
	})
	return vb
}

// WithFirstFrame sets the first frame image (convenience method)
func (vb *VideoRequestBuilder) WithFirstFrame(imageUUID string) *VideoRequestBuilder {
	return vb.WithFrameImage(imageUUID, models.FramePositionFirst)
}

// WithLastFrame sets the last frame image (convenience method)
func (vb *VideoRequestBuilder) WithLastFrame(imageUUID string) *VideoRequestBuilder {
	return vb.WithFrameImage(imageUUID, models.FramePositionLast)
}

// WithSeed sets the seed for reproducibility
func (vb *VideoRequestBuilder) WithSeed(seed int64) *VideoRequestBuilder {
	vb.req.Seed = &seed
	return vb
}

// WithCFGScale sets the CFG scale
func (vb *VideoRequestBuilder) WithCFGScale(scale float64) *VideoRequestBuilder {
	vb.req.CFGScale = &scale
	return vb
}

// WithReferenceImage adds a reference image
func (vb *VideoRequestBuilder) WithReferenceImage(imageUUID string) *VideoRequestBuilder {
	vb.req.ReferenceImages = append(vb.req.ReferenceImages, models.ReferenceImage{
		InputImage: imageUUID,
	})
	return vb
}

// WithReferenceVideo adds a reference video
func (vb *VideoRequestBuilder) WithReferenceVideo(videoUUID string) *VideoRequestBuilder {
	vb.req.ReferenceVideos = append(vb.req.ReferenceVideos, models.ReferenceVideo{
		InputVideo: videoUUID,
	})
	return vb
}

// WithInputAudio adds an input audio
func (vb *VideoRequestBuilder) WithInputAudio(audioUUID string) *VideoRequestBuilder {
	vb.req.InputAudios = append(vb.req.InputAudios, models.InputAudio{
		InputAudio: audioUUID,
	})
	return vb
}

// WithSpeech adds text-to-speech generation
func (vb *VideoRequestBuilder) WithSpeech(voice, text string) *VideoRequestBuilder {
	vb.req.Speech = &models.Speech{
		Voice: voice,
		Text:  text,
	}
	return vb
}

// WithSafety enables content safety checking
func (vb *VideoRequestBuilder) WithSafety(mode models.SafetyMode) *VideoRequestBuilder {
	checkContent := true
	vb.req.Safety = &models.Safety{
		CheckContent: checkContent,
		Mode:         mode,
	}
	return vb
}

// WithLoRA adds a LoRA model
func (vb *VideoRequestBuilder) WithLoRA(model string, weight float64) *VideoRequestBuilder {
	vb.req.LoRA = append(vb.req.LoRA, models.LoRA{
		Model:  model,
		Weight: &weight,
	})
	return vb
}

// WithOutputFormat sets the output format
func (vb *VideoRequestBuilder) WithOutputFormat(format models.VideoOutputFormat) *VideoRequestBuilder {
	vb.req.OutputFormat = &format
	return vb
}

// WithOutputQuality sets the output quality (20-99)
func (vb *VideoRequestBuilder) WithOutputQuality(quality int) *VideoRequestBuilder {
	vb.req.OutputQuality = &quality
	return vb
}

// WithGoogleSettings adds Google (Veo) provider settings
func (vb *VideoRequestBuilder) WithGoogleSettings(enhancePrompt, generateAudio bool) *VideoRequestBuilder {
	if vb.req.ProviderSettings == nil {
		vb.req.ProviderSettings = &models.VideoProviderSettings{}
	}
	vb.req.ProviderSettings.Google = &models.GoogleVideoSettings{
		EnhancePrompt: &enhancePrompt,
		GenerateAudio: &generateAudio,
	}
	return vb
}

// WithPixVerseSettings adds PixVerse provider settings
func (vb *VideoRequestBuilder) WithPixVerseSettings(style, effect, cameraMovement string) *VideoRequestBuilder {
	if vb.req.ProviderSettings == nil {
		vb.req.ProviderSettings = &models.VideoProviderSettings{}
	}
	vb.req.ProviderSettings.PixVerse = &models.PixVerseVideoSettings{
		Style:          &style,
		Effect:         &effect,
		CameraMovement: &cameraMovement,
	}
	return vb
}

// WithViduSettings adds Vidu provider settings
func (vb *VideoRequestBuilder) WithViduSettings(movementAmplitude, style string, bgm bool) *VideoRequestBuilder {
	if vb.req.ProviderSettings == nil {
		vb.req.ProviderSettings = &models.VideoProviderSettings{}
	}
	vb.req.ProviderSettings.Vidu = &models.ViduVideoSettings{
		MovementAmplitude: &movementAmplitude,
		Style:             &style,
		BGM:               &bgm,
	}
	return vb
}

// WithIncludeCost includes cost in the response
func (vb *VideoRequestBuilder) WithIncludeCost(include bool) *VideoRequestBuilder {
	vb.req.IncludeCost = &include
	return vb
}

// Build returns the built video request
func (vb *VideoRequestBuilder) Build() *models.VideoInferenceRequest {
	return vb.req
}

// AudioInference generates audio using the full request object
func (c *Client) AudioInference(
	ctx context.Context,
	request *models.AudioInferenceRequest,
) (*models.AudioInferenceResponse, error) {
	result, err := c.sendRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	if resp, ok := result.(*models.AudioInferenceResponse); ok {
		return resp, nil
	}

	return nil, ErrInvalidResponse
}

// TextToAudio is a convenience method for simple text-to-audio generation
func (c *Client) TextToAudio(
	ctx context.Context,
	prompt string,
	model string,
	duration int,
) (*models.AudioInferenceResponse, error) {
	req := models.NewAudioInferenceRequest(prompt, model, duration)
	return c.AudioInference(ctx, req)
}

// PollAudioResult polls for audio generation result with automatic retries
// This function will keep polling until it gets a definitive response (success/error)
// or reaches maxAttempts. Each individual poll has its own timeout, but the overall
// polling process can run as long as needed.
func (c *Client) PollAudioResult(
	ctx context.Context,
	taskUUID string,
	maxAttempts int,
	pollInterval time.Duration,
) (*models.AudioInferenceResponse, error) {
	c.debugLogger.Printf("Starting audio polling for TaskUUID: %s (max attempts: %d, interval: %v)",
		taskUUID, maxAttempts, pollInterval)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// Check if parent context was canceled
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Create a fresh context for each poll with a reasonable timeout
		// This allows the overall polling to continue beyond the client's default timeout
		pollCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		resp, err := c.GetResponse(pollCtx, taskUUID)
		cancel()

		if err != nil {
			c.debugLogger.Printf("Poll attempt %d/%d: error - %v", attempt, maxAttempts, err)
			// If it's a timeout on this specific poll, continue trying
			if IsTimeout(err) || err == context.DeadlineExceeded {
				c.debugLogger.Printf("Poll request timed out, will retry...")
				time.Sleep(pollInterval)
				continue
			}
			// For other errors, return immediately
			return nil, err
		}

		// The response type depends on what we're polling for
		// Try to cast to AudioInferenceResponse
		if audioResp, ok := resp.(*models.AudioInferenceResponse); ok {
			c.debugLogger.Printf("Poll attempt %d/%d: status = %s", attempt, maxAttempts, audioResp.Status)

			switch audioResp.Status {
			case models.TaskStatusSuccess:
				c.debugLogger.Printf("Audio generation completed successfully!")
				return audioResp, nil
			case models.TaskStatusError:
				return nil, fmt.Errorf("audio generation failed - check API response for details")
			case models.TaskStatusProcessing, "":
				// Continue polling - status is still processing or not set yet
				c.debugLogger.Printf("Audio still processing, waiting %v before next poll...", pollInterval)
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(pollInterval):
					continue
				}
			default:
				// Unknown status, log it but continue polling
				c.debugLogger.Printf("Unknown status: %s, continuing to poll...", audioResp.Status)
				time.Sleep(pollInterval)
			}
		}
	}

	return nil, fmt.Errorf("polling exhausted: reached max attempts (%d) without definitive response", maxAttempts)
}

// AudioRequestBuilder provides a fluent interface for building audio inference requests
type AudioRequestBuilder struct{ req *models.AudioInferenceRequest }

// NewAudioRequestBuilder creates a new audio request builder
func NewAudioRequestBuilder(prompt, model string, duration int) *AudioRequestBuilder {
	return &AudioRequestBuilder{req: models.NewAudioInferenceRequest(prompt, model, duration)}
}

// WithAudioSettings sets the audio quality settings
func (ab *AudioRequestBuilder) WithAudioSettings(sampleRate, bitrate int) *AudioRequestBuilder {
	ab.req.AudioSettings = &models.AudioSettings{
		SampleRate: &sampleRate,
		Bitrate:    &bitrate,
	}
	return ab
}

// WithElevenLabsMusic adds ElevenLabs music generation settings
func (ab *AudioRequestBuilder) WithElevenLabsMusic(promptInfluence float64) *AudioRequestBuilder {
	if ab.req.ProviderSettings == nil {
		ab.req.ProviderSettings = &models.AudioProviderSettings{}
	}
	ab.req.ProviderSettings.ElevenLabs = &models.ElevenLabsAudioSettings{
		Music: &models.ElevenLabsMusicSettings{
			PromptInfluence: &promptInfluence,
		},
	}
	return ab
}

// WithNumberResults sets how many audio files to generate
func (ab *AudioRequestBuilder) WithNumberResults(count int) *AudioRequestBuilder {
	ab.req.NumberResults = &count
	return ab
}

// WithIncludeCost includes cost in the response
func (ab *AudioRequestBuilder) WithIncludeCost(include bool) *AudioRequestBuilder {
	ab.req.IncludeCost = &include
	return ab
}

// WithUploadEndpoint sets where to upload the generated audio
func (ab *AudioRequestBuilder) WithUploadEndpoint(endpoint string) *AudioRequestBuilder {
	ab.req.UploadEndpoint = &endpoint
	return ab
}

// Build returns the built audio request
func (ab *AudioRequestBuilder) Build() *models.AudioInferenceRequest { return ab.req }
