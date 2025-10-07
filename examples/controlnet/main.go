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

	// Upload a guide image for ControlNet
	fmt.Println("Uploading guide image...")
	uploadResp, err := client.UploadImageFromURL(ctx, "https://example.com/pose.jpg")
	if err != nil {
		log.Fatalf("Failed to upload guide image: %v", err)
	}

	fmt.Printf("Guide image uploaded with UUID: %s\n", uploadResp.ImageUUID)

	// Generate image with ControlNet
	request := runware.NewRequestBuilder(
		"a photorealistic portrait of a young woman, professional lighting, high detail",
		"runware:101@1",
		1024,
		1024,
	).
		WithControlNet("runware:25@1", uploadResp.ImageUUID, 0.8).
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
