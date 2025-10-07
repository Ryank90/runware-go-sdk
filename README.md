# Runware Go SDK

A Go SDK for the [Runware AI](https://runware.ai) platform. Generate, transform, and enhance images and videos using state-of-the-art AI models through a simple interface.

## Installation

```bash
go get github.com/Ryank90/runware-go-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    runware "github.com/Ryank90/runware-go-sdk"
)

func main() {
    // Create client (reads RUNWARE_API_KEY from environment)
    client, err := runware.NewClient(nil)
    if err != nil {
        log.Fatal(err)
    }

    // Connect to the API
    ctx := context.Background()
    if err := client.Connect(ctx); err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect()

    // Generate an image
    response, err := client.TextToImage(
        ctx,
        "a serene mountain landscape with a crystal-clear lake",
        "runware:101@1",
        1024,
        1024,
    )
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Image URL: %s\n", *response.ImageURL)
}
```

## Configuration

### Using Environment Variables

```bash
export RUNWARE_API_KEY="your-api-key-here"
```

```go
client, err := runware.NewClient(nil)
```

### Custom Configuration

```go
config := &runware.Config{
    APIKey: "your-api-key-here",
    RequestTimeout: 60 * time.Second,
    WSConfig: &runware.WSConfig{
        URL: runware.DefaultWSURL,
        ConnectTimeout: 30 * time.Second,
        EnableAutoReconnect: true,
    },
}

client, err := runware.NewClient(config)
```

## Usage Examples

See the [`examples/`](./examples) directory for complete, working examples:

**Image Generation:**
- [Text to Image](./examples/text_to_image) - Simple image generation
- [Advanced Generation](./examples/advanced_generation) - Using the builder pattern with advanced options
- [Image Transformation](./examples/image_transformation) - Transform existing images
- [Batch Generation](./examples/batch_generation) - Generate multiple images in parallel
- [ControlNet](./examples/controlnet) - Guided generation with ControlNet
- [Utilities](./examples/utilities) - Upscaling, background removal, captioning, etc.

**Video Generation:**
- [Text to Video](./examples/text_to_video) - Generate videos from text with async polling
- [Image to Video](./examples/image_to_video) - Animate images into videos
- [Advanced Video](./examples/advanced_video) - Provider settings and advanced controls
- [Video with Constraints](./examples/video_with_constraints) - Frame-by-frame constraints
- [Batch Video](./examples/batch_video) - Generate multiple videos in parallel

### Quick Start

```go
// Text to Image
response, err := client.TextToImage(ctx, "a sunset", "runware:101@1", 1024, 1024)

// Text to Video (with async polling)
videoResp, err := client.TextToVideo(ctx, "ocean waves", "klingai:5@3", 5)
finalResp, err := client.PollVideoResult(ctx, videoResp.TaskUUID, 60, 10*time.Second)
```

For detailed usage, see the [examples directory](./examples).

## Models

The SDK supports various model formats:

See the [Image Models documentation](https://runware.ai/docs/en/image-inference/models) for available models.

See the [Video Models documentation](https://runware.ai/docs/en/video-inference/api-reference) for all available models and their capabilities.

## Schedulers

Available schedulers for diffusion sampling:

- `SchedulerEuler`
- `SchedulerEulerA`
- `SchedulerDPMPP2M`
- `SchedulerDPMPP2MKarras`
- `SchedulerDPMPPSDE`
- `SchedulerDPMPPSDEKarras`
- `SchedulerLMS`
- `SchedulerLMSKarras`
- `SchedulerHeun`
- `SchedulerDDIM`
- `SchedulerPNDM`

## API Reference

### Client Methods

#### Connection Management

- `Connect(ctx context.Context) error` - Establish WebSocket connection
- `Disconnect() error` - Close WebSocket connection
- `IsConnected() bool` - Check connection status

#### Image Generation

- `TextToImage(ctx, prompt, model, width, height) (*ImageInferenceResponse, error)`
- `ImageToImage(ctx, prompt, model, seedImage, width, height, strength) (*ImageInferenceResponse, error)`
- `Inpaint(ctx, prompt, model, seedImage, maskImage, width, height, strength) (*ImageInferenceResponse, error)`
- `Outpaint(ctx, prompt, model, seedImage, width, height, outpaint) (*ImageInferenceResponse, error)`
- `ImageInference(ctx, request) (*ImageInferenceResponse, error)`
- `ImageInferenceBatch(ctx, requests) ([]*ImageInferenceResponse, error)`

#### Video Generation

- `TextToVideo(ctx, prompt, model, duration) (*VideoInferenceResponse, error)`
- `ImageToVideo(ctx, prompt, model, seedImage, duration) (*VideoInferenceResponse, error)`
- `VideoInference(ctx, request) (*VideoInferenceResponse, error)`
- `VideoInferenceBatch(ctx, requests) ([]*VideoInferenceResponse, error)`

#### Image Utilities

- `UploadImage(ctx, request) (*UploadImageResponse, error)`
- `UploadImageFromFile(ctx, filePath) (*UploadImageResponse, error)`
- `UploadImageFromURL(ctx, url) (*UploadImageResponse, error)`
- `UpscaleImage(ctx, request) (*UpscaleGanResponse, error)`
- `RemoveBackground(ctx, request) (*RemoveImageBackgroundResponse, error)`

#### Text Utilities

- `EnhancePrompt(ctx, request) (*EnhancePromptResponse, error)`
- `CaptionImage(ctx, request) (*ImageCaptionResponse, error)`


## Testing

Run the test suite:

```bash
go test -v ./...
```

Run tests with coverage:

```bash
go test -v -cover ./...
```

## Support

- **Documentation**: [https://runware.ai/docs](https://runware.ai/docs)
- **API Reference**: [https://runware.ai/docs/en/image-inference/api-reference](https://runware.ai/docs/en/image-inference/api-reference)
- **GitHub Issues**: [https://github.com/Ryank90/runware-go-sdk/issues](https://github.com/Ryank90/runware-go-sdk/issues)