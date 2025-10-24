// Package models provides type definitions for the Runware AI platform API.
//
// This package contains all request and response types, enums, and constants
// used to interact with the Runware API. Types are organized by domain:
//
//   - image_types.go: Image generation, upload, upscaling, background removal
//   - video_types.go: Video generation and configuration
//   - audio_types.go: Audio/music generation
//   - shared_types.go: Common types, enums, and constants
//   - constructors.go: Helper functions to create properly initialized requests
//
// # Request Constructors
//
// Each request type has a corresponding New* constructor function that creates
// a properly initialized request with required fields:
//
//	// Image generation
//	req := models.NewImageInferenceRequest(prompt, model, width, height)
//
//	// Video generation
//	req := models.NewVideoInferenceRequest(prompt, model)
//
//	// Audio generation
//	req := models.NewAudioInferenceRequest(prompt, model, duration)
//
// # Constants and Enums
//
// The package exports numerous constants for configuration:
//
//	// Output formats
//	models.OutputTypeURL
//	models.OutputTypeBase64Data
//	models.OutputFormatPNG
//	models.OutputFormatJPG
//
//	// Schedulers
//	models.SchedulerDPMPP2M
//	models.SchedulerEulerA
//
//	// Task types
//	models.TaskTypeImageInference
//	models.TaskTypeVideoInference
//
// # Optional Fields
//
// Most request fields are optional (pointers). Set only what you need:
//
//	req := models.NewImageInferenceRequest("sunset", "runware:101@1", 1024, 1024)
//	steps := 30
//	req.Steps = &steps  // Optional: customize step count
//	cfg := 7.5
//	req.CFGScale = &cfg  // Optional: customize CFG scale
//
// # Provider-Specific Settings
//
// Video and audio requests support provider-specific settings:
//
//	req := models.NewVideoInferenceRequest("ocean waves", "klingai:5@3")
//	req.ProviderSettings = &models.VideoProviderSettings{
//	    KlingAI: &models.KlingAIVideoSettings{
//	        Mode: stringPtr("professional"),
//	    },
//	}
//
// # Response Types
//
// Response types contain both required and optional fields. Use pointer
// checks before dereferencing optional fields:
//
//	resp, err := client.TextToImage(ctx, prompt, model, width, height)
//	if err != nil {
//	    return err
//	}
//
//	// Required field
//	fmt.Println("Image UUID:", resp.ImageUUID)
//
//	// Optional field - check before use
//	if resp.ImageURL != nil {
//	    fmt.Println("Image URL:", *resp.ImageURL)
//	}
//	if resp.Cost != nil {
//	    fmt.Printf("Cost: $%.4f\n", *resp.Cost)
//	}
//
// # Advanced Features
//
// The package supports advanced AI features:
//
//	// ControlNet for guided generation
//	controlNet := models.ControlNet{
//	    Model:      "canny",
//	    GuideImage: imageUUID,
//	}
//	req.ControlNet = []models.ControlNet{controlNet}
//
//	// LoRA models
//	lora := models.LoRA{
//	    Model: "style-lora",
//	    Weight: floatPtr(0.8),
//	}
//	req.LoRA = []models.LoRA{lora}
//
//	// IP Adapters
//	adapter := models.IPAdapter{
//	    Model:  "ip-adapter-plus",
//	    Images: []string{imageUUID},
//	}
//	req.IPAdapters = []models.IPAdapter{adapter}
package models
