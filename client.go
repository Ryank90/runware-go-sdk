package runware

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"sync"
	"time"
)

// Client is the main Runware SDK client
type Client struct {
	ws             *wsClient
	apiKey         string
	config         *Config
	requestTimeout time.Duration
}

// Config contains client configuration options
type Config struct {
	// APIKey is the Runware API key
	APIKey string

	// WebSocket configuration
	WSConfig *WSConfig

	// RequestTimeout is the default timeout for API requests
	RequestTimeout time.Duration
}

// DefaultConfig returns a default client configuration
func DefaultConfig(apiKey string) *Config {
	return &Config{
		APIKey:         apiKey,
		WSConfig:       DefaultWSConfig(),
		RequestTimeout: 60 * time.Second,
	}
}

// NewClient creates a new Runware client
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		apiKey := os.Getenv("RUNWARE_API_KEY")
		if apiKey == "" {
			return nil, ErrInvalidAPIKey
		}
		config = DefaultConfig(apiKey)
	}

	if config.APIKey == "" {
		return nil, ErrInvalidAPIKey
	}

	client := &Client{
		apiKey:         config.APIKey,
		config:         config,
		requestTimeout: config.RequestTimeout,
		ws:             newWSClient(config.APIKey, config.WSConfig),
	}

	return client, nil
}

// Connect establishes a connection to the Runware API
func (c *Client) Connect(ctx context.Context) error {
	return c.ws.Connect(ctx)
}

// Disconnect closes the connection to the Runware API
func (c *Client) Disconnect() error {
	return c.ws.Disconnect()
}

// IsConnected returns whether the client is connected
func (c *Client) IsConnected() bool {
	return c.ws.IsConnected()
}

// ImageInference performs image inference
func (c *Client) ImageInference(ctx context.Context, req *ImageInferenceRequest) (*ImageInferenceResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	// Set default task type if not set
	if req.TaskType == "" {
		req.TaskType = TaskTypeImageInference
	}

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*ImageInferenceResponse), nil
}

// batchResult is a generic structure for batch processing results
type batchResult[T any] struct {
	index int
	resp  T
	err   error
}

// processBatch processes multiple requests in parallel using a generic handler
func processBatch[Req any, Resp any](
	ctx context.Context,
	requests []Req,
	handler func(context.Context, Req) (Resp, error),
) ([]Resp, error) {
	if len(requests) == 0 {
		return nil, ErrInvalidRequest
	}

	results := make(chan batchResult[Resp], len(requests))
	var wg sync.WaitGroup

	for i, req := range requests {
		wg.Add(1)
		go func(idx int, request Req) {
			defer wg.Done()
			resp, err := handler(ctx, request)
			results <- batchResult[Resp]{index: idx, resp: resp, err: err}
		}(i, req)
	}

	wg.Wait()
	close(results)

	responses := make([]Resp, len(requests))
	var errs []error

	for r := range results {
		if r.err != nil {
			errs = append(errs, fmt.Errorf("request %d: %w", r.index, r.err))
		} else {
			responses[r.index] = r.resp
		}
	}

	if len(errs) > 0 {
		return responses, fmt.Errorf("batch errors: %v", errs)
	}

	return responses, nil
}

// ImageInferenceBatch performs multiple image inference requests in parallel
func (c *Client) ImageInferenceBatch(ctx context.Context, requests []*ImageInferenceRequest) ([]*ImageInferenceResponse, error) {
	return processBatch(ctx, requests, c.ImageInference)
}

// UploadImage uploads an image to Runware
func (c *Client) UploadImage(ctx context.Context, req *UploadImageRequest) (*UploadImageResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.TaskType == "" {
		req.TaskType = TaskTypeImageUpload
	}

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*UploadImageResponse), nil
}

// UploadImageFromFile uploads an image from a file path
func (c *Client) UploadImageFromFile(ctx context.Context, filePath string) (*UploadImageResponse, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	base64Data := base64.StdEncoding.EncodeToString(data)
	req := NewUploadImageRequest()
	req.ImageBase64 = &base64Data

	return c.UploadImage(ctx, req)
}

// UploadImageFromURL uploads an image from a URL
func (c *Client) UploadImageFromURL(ctx context.Context, url string) (*UploadImageResponse, error) {
	req := NewUploadImageRequest()
	req.ImageURL = &url

	return c.UploadImage(ctx, req)
}

// UpscaleImage upscales an image using GAN
func (c *Client) UpscaleImage(ctx context.Context, req *UpscaleGanRequest) (*UpscaleGanResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.TaskType == "" {
		req.TaskType = TaskTypeUpscaleGan
	}

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*UpscaleGanResponse), nil
}

// RemoveBackground removes the background from an image
func (c *Client) RemoveBackground(ctx context.Context, req *RemoveImageBackgroundRequest) (*RemoveImageBackgroundResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.TaskType == "" {
		req.TaskType = TaskTypeImageBackgroundRemoval
	}

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*RemoveImageBackgroundResponse), nil
}

// EnhancePrompt enhances a text prompt
func (c *Client) EnhancePrompt(ctx context.Context, req *EnhancePromptRequest) (*EnhancePromptResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.TaskType == "" {
		req.TaskType = TaskTypePromptEnhance
	}

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*EnhancePromptResponse), nil
}

// CaptionImage generates a caption for an image
func (c *Client) CaptionImage(ctx context.Context, req *ImageCaptionRequest) (*ImageCaptionResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.TaskType == "" {
		req.TaskType = TaskTypeImageCaption
	}

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*ImageCaptionResponse), nil
}

// VideoInference performs video inference
func (c *Client) VideoInference(ctx context.Context, req *VideoInferenceRequest) (*VideoInferenceResponse, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}

	if req.TaskType == "" {
		req.TaskType = TaskTypeVideoInference
	}

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*VideoInferenceResponse), nil
}

// VideoInferenceBatch performs multiple video inference requests in parallel
func (c *Client) VideoInferenceBatch(ctx context.Context, requests []*VideoInferenceRequest) ([]*VideoInferenceResponse, error) {
	return processBatch(ctx, requests, c.VideoInference)
}

// sendRequest is a generic method to send a request and wait for response
func (c *Client) sendRequest(ctx context.Context, req interface{}) (interface{}, error) {
	if !c.IsConnected() {
		return nil, ErrNotConnected
	}

	// Create a channel to receive the response
	respChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)

	handler := func(data interface{}, err error) {
		if err != nil {
			errChan <- err
		} else {
			respChan <- data
		}
	}

	// Send the request
	if err := c.ws.Send(ctx, req, handler); err != nil {
		return nil, err
	}

	// Wait for response with timeout
	timeout := c.requestTimeout
	if deadline, ok := ctx.Deadline(); ok {
		timeout = time.Until(deadline)
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errChan:
		return nil, err
	case resp := <-respChan:
		return resp, nil
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}

// Helper methods for common operations

// TextToImage generates an image from a text prompt
func (c *Client) TextToImage(ctx context.Context, prompt, model string, width, height int) (*ImageInferenceResponse, error) {
	req := NewImageInferenceRequest(prompt, model, width, height)
	return c.ImageInference(ctx, req)
}

// TextToVideo generates a video from a text prompt
func (c *Client) TextToVideo(ctx context.Context, prompt, model string, duration int) (*VideoInferenceResponse, error) {
	req := NewVideoInferenceRequest(prompt, model)
	req.Duration = &duration
	return c.VideoInference(ctx, req)
}

// ImageToVideo generates a video from an image and prompt
func (c *Client) ImageToVideo(ctx context.Context, prompt, model, seedImage string, duration int) (*VideoInferenceResponse, error) {
	req := NewVideoInferenceRequest(prompt, model)
	req.Duration = &duration
	req.FrameImages = []FrameImage{{InputImage: seedImage, Frame: FramePositionFirst}}
	return c.VideoInference(ctx, req)
}

// GetResponse polls for the result of an async task
func (c *Client) GetResponse(ctx context.Context, taskUUID string) (*VideoInferenceResponse, error) {
	req := NewGetResponseRequest(taskUUID)

	result, err := c.sendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	return result.(*VideoInferenceResponse), nil
}

// PollVideoResult polls for video generation result with automatic retries
func (c *Client) PollVideoResult(
	ctx context.Context,
	taskUUID string,
	maxAttempts int,
	pollInterval time.Duration,
) (*VideoInferenceResponse, error) {
	for attempt := 0; attempt < maxAttempts; attempt++ {
		resp, err := c.GetResponse(ctx, taskUUID)
		if err != nil {
			return nil, err
		}

		switch resp.Status {
		case TaskStatusSuccess:
			return resp, nil
		case TaskStatusError:
			return nil, fmt.Errorf("video generation failed")
		case TaskStatusProcessing:
			// Continue polling
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(pollInterval):
				continue
			}
		}
	}

	return nil, fmt.Errorf("polling timeout: max attempts (%d) reached", maxAttempts)
}

// ImageToImage transforms an image based on a prompt
func (c *Client) ImageToImage(
	ctx context.Context,
	prompt, model, seedImage string,
	width, height int,
	strength float64,
) (*ImageInferenceResponse, error) {
	req := NewImageInferenceRequest(prompt, model, width, height)
	req.SeedImage = &seedImage
	req.Strength = &strength
	return c.ImageInference(ctx, req)
}

// Inpaint performs inpainting on an image
func (c *Client) Inpaint(
	ctx context.Context,
	prompt, model, seedImage, maskImage string,
	width, height int,
	strength float64,
) (*ImageInferenceResponse, error) {
	req := NewImageInferenceRequest(prompt, model, width, height)
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
	outpaint *Outpaint,
) (*ImageInferenceResponse, error) {
	req := NewImageInferenceRequest(prompt, model, width, height)
	req.SeedImage = &seedImage
	req.Outpaint = outpaint
	return c.ImageInference(ctx, req)
}

// RequestBuilder provides a fluent interface for building requests
type RequestBuilder struct {
	req *ImageInferenceRequest
}

// NewRequestBuilder creates a new request builder
func NewRequestBuilder(prompt, model string, width, height int) *RequestBuilder {
	return &RequestBuilder{
		req: NewImageInferenceRequest(prompt, model, width, height),
	}
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
func (rb *RequestBuilder) WithScheduler(scheduler Scheduler) *RequestBuilder {
	rb.req.Scheduler = &scheduler
	return rb
}

// WithNumberResults sets the number of results
func (rb *RequestBuilder) WithNumberResults(num int) *RequestBuilder {
	rb.req.NumberResults = &num
	return rb
}

// WithOutputType sets the output type
func (rb *RequestBuilder) WithOutputType(outputType OutputType) *RequestBuilder {
	rb.req.OutputType = &outputType
	return rb
}

// WithOutputFormat sets the output format
func (rb *RequestBuilder) WithOutputFormat(format OutputFormat) *RequestBuilder {
	rb.req.OutputFormat = &format
	return rb
}

// WithLoRA adds a LoRA
func (rb *RequestBuilder) WithLoRA(model string, weight float64) *RequestBuilder {
	rb.req.LoRA = append(rb.req.LoRA, LoRA{
		Model:  model,
		Weight: &weight,
	})
	return rb
}

// WithControlNet adds a ControlNet
func (rb *RequestBuilder) WithControlNet(model, guideImage string, weight float64) *RequestBuilder {
	rb.req.ControlNet = append(rb.req.ControlNet, ControlNet{
		Model:      model,
		GuideImage: guideImage,
		Weight:     &weight,
	})
	return rb
}

// WithEmbedding adds an embedding
func (rb *RequestBuilder) WithEmbedding(model string, weight float64) *RequestBuilder {
	rb.req.Embeddings = append(rb.req.Embeddings, Embedding{
		Model:  model,
		Weight: &weight,
	})
	return rb
}

// WithIPAdapter adds an IP-Adapter
func (rb *RequestBuilder) WithIPAdapter(model, guideImage string, weight float64) *RequestBuilder {
	rb.req.IPAdapters = append(rb.req.IPAdapters, IPAdapter{
		Model:      model,
		GuideImage: guideImage,
		Weight:     &weight,
	})
	return rb
}

// WithOutpaint sets the outpaint parameters
func (rb *RequestBuilder) WithOutpaint(outpaint *Outpaint) *RequestBuilder {
	rb.req.Outpaint = outpaint
	return rb
}

// WithRefiner sets the refiner
func (rb *RequestBuilder) WithRefiner(model string, startStep int) *RequestBuilder {
	rb.req.Refiner = &Refiner{
		Model:     model,
		StartStep: &startStep,
	}
	return rb
}

// WithSafety enables safety checks
func (rb *RequestBuilder) WithSafety(mode SafetyMode) *RequestBuilder {
	checkContent := true
	rb.req.Safety = &Safety{
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
func (rb *RequestBuilder) Build() *ImageInferenceRequest {
	return rb.req
}

// VideoRequestBuilder provides a fluent interface for building video requests
type VideoRequestBuilder struct {
	req *VideoInferenceRequest
}

// NewVideoRequestBuilder creates a new video request builder
func NewVideoRequestBuilder(prompt, model string) *VideoRequestBuilder {
	return &VideoRequestBuilder{
		req: NewVideoInferenceRequest(prompt, model),
	}
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
func (vb *VideoRequestBuilder) WithFrameImage(imageUUID string, position FramePosition) *VideoRequestBuilder {
	vb.req.FrameImages = append(vb.req.FrameImages, FrameImage{
		InputImage: imageUUID,
		Frame:      position,
	})
	return vb
}

// WithFirstFrame sets the first frame image (convenience method)
func (vb *VideoRequestBuilder) WithFirstFrame(imageUUID string) *VideoRequestBuilder {
	return vb.WithFrameImage(imageUUID, FramePositionFirst)
}

// WithLastFrame sets the last frame image (convenience method)
func (vb *VideoRequestBuilder) WithLastFrame(imageUUID string) *VideoRequestBuilder {
	return vb.WithFrameImage(imageUUID, FramePositionLast)
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
	vb.req.ReferenceImages = append(vb.req.ReferenceImages, ReferenceImage{
		InputImage: imageUUID,
	})
	return vb
}

// WithReferenceVideo adds a reference video
func (vb *VideoRequestBuilder) WithReferenceVideo(videoUUID string) *VideoRequestBuilder {
	vb.req.ReferenceVideos = append(vb.req.ReferenceVideos, ReferenceVideo{
		InputVideo: videoUUID,
	})
	return vb
}

// WithInputAudio adds an input audio
func (vb *VideoRequestBuilder) WithInputAudio(audioUUID string) *VideoRequestBuilder {
	vb.req.InputAudios = append(vb.req.InputAudios, InputAudio{
		InputAudio: audioUUID,
	})
	return vb
}

// WithSpeech adds text-to-speech generation
func (vb *VideoRequestBuilder) WithSpeech(voice, text string) *VideoRequestBuilder {
	vb.req.Speech = &Speech{
		Voice: voice,
		Text:  text,
	}
	return vb
}

// WithSafety enables content safety checking
func (vb *VideoRequestBuilder) WithSafety(mode SafetyMode) *VideoRequestBuilder {
	checkContent := true
	vb.req.Safety = &Safety{
		CheckContent: checkContent,
		Mode:         mode,
	}
	return vb
}

// WithLoRA adds a LoRA model
func (vb *VideoRequestBuilder) WithLoRA(model string, weight float64) *VideoRequestBuilder {
	vb.req.LoRA = append(vb.req.LoRA, LoRA{
		Model:  model,
		Weight: &weight,
	})
	return vb
}

// WithOutputFormat sets the output format
func (vb *VideoRequestBuilder) WithOutputFormat(format VideoOutputFormat) *VideoRequestBuilder {
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
		vb.req.ProviderSettings = &VideoProviderSettings{}
	}
	vb.req.ProviderSettings.Google = &GoogleVideoSettings{
		EnhancePrompt: &enhancePrompt,
		GenerateAudio: &generateAudio,
	}
	return vb
}

// WithPixVerseSettings adds PixVerse provider settings
func (vb *VideoRequestBuilder) WithPixVerseSettings(style, effect, cameraMovement string) *VideoRequestBuilder {
	if vb.req.ProviderSettings == nil {
		vb.req.ProviderSettings = &VideoProviderSettings{}
	}
	vb.req.ProviderSettings.PixVerse = &PixVerseVideoSettings{
		Style:          &style,
		Effect:         &effect,
		CameraMovement: &cameraMovement,
	}
	return vb
}

// WithViduSettings adds Vidu provider settings
func (vb *VideoRequestBuilder) WithViduSettings(movementAmplitude, style string, bgm bool) *VideoRequestBuilder {
	if vb.req.ProviderSettings == nil {
		vb.req.ProviderSettings = &VideoProviderSettings{}
	}
	vb.req.ProviderSettings.Vidu = &ViduVideoSettings{
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
func (vb *VideoRequestBuilder) Build() *VideoInferenceRequest {
	return vb.req
}
