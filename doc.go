// Package runware provides a comprehensive Go SDK for the Runware AI platform.
//
// # Overview
//
// The Runware SDK enables developers to integrate state-of-the-art AI capabilities
// into their Go applications, including:
//
//   - Image Generation: Text-to-image, image-to-image, inpainting, outpainting
//   - Video Generation: Text-to-video, image-to-video with multiple provider support
//   - Audio Generation: Text-to-audio/music with high-quality synthesis
//   - Image Utilities: Upscaling, background removal, prompt enhancement, captioning
//
// # Quick Start
//
// Create a client and generate an image:
//
//	package main
//
//	import (
//	    "context"
//	    "fmt"
//	    "log"
//
//	    "github.com/Ryank90/runware-go-sdk"
//	)
//
//	func main() {
//	    // Create client (reads RUNWARE_API_KEY from environment)
//	    client, err := runware.NewClient(nil)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    defer client.Disconnect()
//
//	    // Connect to the API
//	    ctx := context.Background()
//	    if err := client.Connect(ctx); err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Generate an image
//	    resp, err := client.TextToImage(ctx,
//	        "a serene mountain landscape with crystal-clear lake",
//	        "runware:101@1", 1024, 1024)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    fmt.Printf("Image URL: %s\n", *resp.ImageURL)
//	}
//
// # Configuration
//
// The SDK can be configured explicitly or via environment variables:
//
//	// Use environment variable RUNWARE_API_KEY
//	client, err := runware.NewClient(nil)
//
//	// Explicit configuration
//	config := runware.DefaultConfig()
//	config.APIKey = "your-api-key"
//	config.RequestTimeout = 60 * time.Second
//	config.EnableDebugLogging = true
//	client, err := runware.NewClient(config)
//
// # Error Handling
//
// The SDK provides detailed error types for robust error handling:
//
//	resp, err := client.TextToImage(ctx, prompt, model, width, height)
//	if err != nil {
//	    switch {
//	    case errors.Is(err, runware.ErrNotConnected):
//	        // Handle connection error
//	    case errors.Is(err, runware.ErrTimeout):
//	        // Handle timeout
//	    default:
//	        if apiErr, ok := err.(*runware.APIError); ok {
//	            fmt.Printf("API Error: %s (ID: %s)\n", apiErr.Message, apiErr.ErrorID)
//	            if apiErr.IsRetryable() {
//	                // Implement retry logic
//	            }
//	        }
//	    }
//	}
//
// # Image Generation
//
// ## Simple Text-to-Image
//
// Use TextToImage for the simplest case:
//
//	resp, err := client.TextToImage(ctx, "sunset over ocean", "runware:101@1", 1024, 1024)
//
// ## Advanced Image Generation
//
// For fine-grained control, use ImageInference with a configured request:
//
//	req := models.NewImageInferenceRequest("mountain landscape", "runware:101@1", 1024, 1024)
//	steps := 30
//	cfg := 7.5
//	req.Steps = &steps
//	req.CFGScale = &cfg
//	req.NegativePrompt = stringPtr("blurry, low quality")
//
//	resp, err := client.ImageInference(ctx, req)
//
// ## Batch Processing
//
// Generate multiple images efficiently with bounded concurrency:
//
//	requests := []*models.ImageInferenceRequest{
//	    models.NewImageInferenceRequest("sunset", "runware:101@1", 1024, 1024),
//	    models.NewImageInferenceRequest("forest", "runware:101@1", 1024, 1024),
//	    models.NewImageInferenceRequest("ocean", "runware:101@1", 1024, 1024),
//	}
//	responses, err := client.ImageInferenceBatch(ctx, requests)
//
// # Video Generation
//
// Video generation is asynchronous. Submit a request, then poll for results:
//
//	// Submit video generation
//	resp, err := client.TextToVideo(ctx, "ocean waves at sunset", "klingai:5@3", 5)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Poll for completion (max 120 attempts, 15 second intervals)
//	final, err := client.PollVideoResult(ctx, resp.TaskUUID, 120, 15*time.Second)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("Video URL: %s\n", *final.VideoURL)
//
// # Audio Generation
//
// Similar to video, audio generation is asynchronous:
//
//	// Submit audio generation
//	resp, err := client.TextToAudio(ctx, "gentle piano melody", "elevenlabs:1@1", 30)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Poll for completion
//	final, err := client.PollAudioResult(ctx, resp.TaskUUID, 60, 5*time.Second)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("Audio URL: %s\n", *final.AudioURL)
//
// # Concurrency
//
// The Client is safe for concurrent use by multiple goroutines. A single client
// instance can efficiently handle multiple simultaneous requests through
// multiplexed WebSocket communication.
//
//	var wg sync.WaitGroup
//	for i := 0; i < 10; i++ {
//	    wg.Add(1)
//	    go func(idx int) {
//	        defer wg.Done()
//	        resp, err := client.TextToImage(ctx, fmt.Sprintf("image %d", idx),
//	            "runware:101@1", 1024, 1024)
//	        // Handle response...
//	    }(i)
//	}
//	wg.Wait()
//
// # Debug Logging
//
// Enable debug logging to troubleshoot connection or API issues:
//
//	// Via environment variable
//	export RUNWARE_DEBUG=1
//
//	// Or programmatically
//	config := runware.DefaultConfig()
//	config.EnableDebugLogging = true
//	client, _ := runware.NewClient(config)
//
// # Models Package
//
// The models package contains all request/response types and constants:
//
//	import "github.com/Ryank90/runware-go-sdk/models"
//
//	// Create requests
//	req := models.NewImageInferenceRequest(prompt, model, width, height)
//
//	// Use constants
//	req.Scheduler = &models.SchedulerDPMPP2M
//	req.OutputFormat = &models.OutputFormatPNG
//
// # Examples
//
// See the examples/ directory in the repository for complete, working examples:
//
//   - examples/text_to_image - Simple image generation
//   - examples/advanced_generation - Advanced options with builder pattern
//   - examples/batch_generation - Parallel batch processing
//   - examples/text_to_video - Video generation with polling
//   - examples/text_to_audio - Audio generation
//   - examples/utilities - Upscaling, background removal, etc.
//
// # Additional Resources
//
//   - Documentation: https://runware.ai/docs
//   - API Reference: https://runware.ai/docs/en/image-inference/api-reference
//   - GitHub: https://github.com/Ryank90/runware-go-sdk
package runware
