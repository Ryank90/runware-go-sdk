package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	runware "github.com/Ryank90/runware-go-sdk"
)

func main() {
	// Create a new client with custom configuration
	config := runware.DefaultConfig()
	config.APIKey = os.Getenv("RUNWARE_API_KEY") // or set directly: "your-api-key-here"
	config.RequestTimeout = 90 * time.Second     // Longer timeout for complex requests
	client, err := runware.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Connect to the API
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect()

	fmt.Println("Connected to Runware API")

	// Build a request with advanced features
	// Using single-result generation for reliability with advanced settings
	outputType := runware.OutputTypeURL
	outputFormat := runware.OutputFormatPNG

	request := runware.NewRequestBuilder(
		"a highly detailed portrait of a wise old wizard with a long beard, fantasy art, dramatic lighting",
		"runware:101@1",
		1024,
		1024,
	).
		WithNegativePrompt("blurry, low quality, distorted").
		WithSteps(40).
		WithCFGScale(7.5).
		WithOutputType(outputType).
		WithOutputFormat(outputFormat).
		WithIncludeCost(true).
		Build()

	fmt.Println("Generating image with advanced settings...")
	fmt.Println("(Using custom steps, CFG scale, negative prompt, and PNG output)")

	response, err := client.ImageInference(ctx, request)
	if err != nil {
		log.Fatalf("Failed to generate image: %v", err)
	}

	fmt.Printf("\nImage generated successfully!\n")
	fmt.Printf("Image UUID: %s\n", response.ImageUUID)
	if response.ImageURL != nil {
		fmt.Printf("Image URL: %s\n", *response.ImageURL)
	}
	if response.Seed != nil {
		fmt.Printf("Seed: %d (use this to reproduce the image)\n", *response.Seed)
	}
	if response.Cost != nil {
		fmt.Printf("Cost: $%.4f\n", *response.Cost)
	}
}
