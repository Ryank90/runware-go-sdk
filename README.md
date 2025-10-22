# Runware Go SDK

A Go SDK for the [Runware AI](https://runware.ai) platform. Generate, transform, and enhance images, videos, and audio using state-of-the-art AI models through a simple interface.

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
    
    "github.com/Ryank90/runware-go-sdk"
    "github.com/Ryank90/runware-go-sdk/models"
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
import (
    "github.com/Ryank90/runware-go-sdk"
    wsinternal "github.com/Ryank90/runware-go-sdk/internal/ws"
)

config := &runware.Config{
    APIKey: "your-api-key-here",
    RequestTimeout: 60 * time.Second,
    WSConfig: wsinternal.DefaultWSConfig(),
}

client, err := runware.NewClient(config)
```

## Error Handling & Debugging

The SDK provides comprehensive error handling with detailed context for production debugging.

### Enhanced Error Messages

```go
resp, err := client.TextToImage(ctx, prompt, model, width, height)
if err != nil {
    // Detailed error with full API context
    // Example: "Runware API Error - Task: imageInference | Code: unsupportedDimensions | Message: ... | TaskUUID: ..."
    fmt.Printf("Error: %v\n", err)
    
    // Check specific error types
    if apiErr, ok := err.(*runware.APIError); ok {
        if apiErr.IsRetryable() {
            // Retry logic for transient errors
        }
    }
}
```

### Debug Logging

Enable detailed logging to troubleshoot issues:

```bash
# Enable debug mode
export RUNWARE_DEBUG=1
```

Or programmatically:

```go
config := runware.DefaultConfig()
config.EnableDebugLogging = true
client, _ := runware.NewClient(config)
```

## Usage Examples

See the [`examples/`](./examples) directory for complete, working examples.

## Models

The SDK supports various model formats:

- **Image Models**: See the [Image Models documentation](https://runware.ai/docs/en/image-inference/models)
- **Video Models**: See the [Video Models documentation](https://runware.ai/docs/en/video-inference/api-reference)
- **Audio Models**: See the [Audio Inference documentation](https://runware.ai/docs/en/audio-inference/api-reference)

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
- `PollVideoResult(ctx, taskUUID, maxAttempts, pollInterval) (*VideoInferenceResponse, error)`

#### Audio Generation

- `TextToAudio(ctx, prompt, model, duration) (*AudioInferenceResponse, error)`
- `AudioInference(ctx, request) (*AudioInferenceResponse, error)`
- `PollAudioResult(ctx, taskUUID, maxAttempts, pollInterval) (*AudioInferenceResponse, error)`

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
- **GitHub Issues**: [https://github.com/Ryank90/runware-go-sdk/issues](https://github.com/Ryank90/runware-go-sdk/issues)