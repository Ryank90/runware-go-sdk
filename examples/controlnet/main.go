package main

import (
	"context"
	"fmt"
	"log"

	runware "github.com/Ryank90/runware-go-sdk"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	client, err := runware.NewClient(nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect()

	fmt.Println("Connected to Runware API")

	// Generate a base image to use as ControlNet guide
	fmt.Println("Generating guide image...")
	guidePrompt := "a simple stick figure pose, line art, white background"
	guideResp, err := client.TextToImage(ctx, guidePrompt, "runware:101@1", 512, 512)
	if err != nil {
		log.Fatalf("Failed to generate guide image: %v", err)
	}

	fmt.Printf("Guide image generated with UUID: %s\n", guideResp.ImageUUID)

	// Get image URL for ControlNet (guide image requires URL)
	if guideResp.ImageURL == nil {
		log.Fatalf("No image URL returned")
	}

	// Generate image with ControlNet
	request := runware.NewRequestBuilder(
		"a photorealistic portrait of a young woman, professional lighting, high detail",
		"runware:101@1",
		512,
		512,
	).
		WithControlNet("runware:25@1", *guideResp.ImageURL, 0.8).
		WithSteps(40).
		WithCFGScale(7.5).
		WithIncludeCost(true).
		Build()

	fmt.Println("Generating image with ControlNet guidance...")

	response, err := client.ImageInference(ctx, request)
	if err != nil {
		log.Fatalf("Failed to generate image: %v", err)
	}

	fmt.Printf("Image generated successfully!\n")
	fmt.Printf("Image UUID: %s\n", response.ImageUUID)
	if response.ImageURL != nil {
		fmt.Printf("Image URL: %s\n", *response.ImageURL)
	}
	if response.Cost != nil {
		fmt.Printf("Cost: $%.4f\n", *response.Cost)
	}
}
