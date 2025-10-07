package main

import (
	"context"
	"fmt"
	"log"

	runware "github.com/Ryank90/runware-go-sdk"
)

func main() {
	// Create a new client with custom configuration
	config := runware.DefaultConfig("your-api-key-here")
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

	// Build a complex request using the fluent builder
	outputType := runware.OutputTypeURL
	outputFormat := runware.OutputFormatPNG
	scheduler := runware.SchedulerDPMPP2MKarras

	request := runware.NewRequestBuilder(
		"a highly detailed portrait of a wise old wizard with a long beard, fantasy art, dramatic lighting",
		"runware:101@1",
		1024,
		1024,
	).
		WithNegativePrompt("blurry, low quality, distorted").
		WithSteps(50).
		WithCFGScale(7.5).
		WithScheduler(scheduler).
		WithNumberResults(4).
		WithOutputType(outputType).
		WithOutputFormat(outputFormat).
		WithIncludeCost(true).
		WithSafety(runware.SafetyModeModerate).
		Build()

	fmt.Println("Generating images with advanced settings...")

	response, err := client.ImageInference(ctx, request)
	if err != nil {
		log.Fatalf("Failed to generate image: %v", err)
	}

	fmt.Printf("Image generated successfully!\n")
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
