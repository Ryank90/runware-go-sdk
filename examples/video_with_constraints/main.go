package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	// Step 1: Generate start frame
	fmt.Println("\nStep 1: Generating start frame...")
	startPrompt := "A peaceful sunrise over misty mountains, photorealistic"
	startFrame, err := client.TextToImage(ctx, startPrompt, "runware:101@1", 1024, 576)
	if err != nil {
		log.Fatalf("Failed to generate start frame: %v", err)
	}
	if startFrame.ImageURL == nil {
		log.Fatalf("No image URL returned for start frame")
	}
	fmt.Printf("Start frame generated: %s\n", startFrame.ImageUUID)
	fmt.Printf("Start frame URL: %s\n", *startFrame.ImageURL)

	// Step 2: Generate end frame
	fmt.Println("\nStep 2: Generating end frame...")
	endPrompt := "The same mountains at sunset with golden light, photorealistic"
	endFrame, err := client.TextToImage(ctx, endPrompt, "runware:101@1", 1024, 576)
	if err != nil {
		log.Fatalf("Failed to generate end frame: %v", err)
	}
	if endFrame.ImageURL == nil {
		log.Fatalf("No image URL returned for end frame")
	}
	fmt.Printf("End frame generated: %s\n", endFrame.ImageUUID)
	fmt.Printf("End frame URL: %s\n", *endFrame.ImageURL)

	// Step 3: Create video with frame constraints for smooth transition
	fmt.Println("\nStep 3: Creating video with frame constraints...")
	request := runware.NewVideoRequestBuilder(
		"Smooth timelapse transition from sunrise to sunset, natural lighting change",
		"klingai:5@3",
	).
		WithDuration(5).
		WithFirstFrame(*startFrame.ImageURL).
		WithLastFrame(*endFrame.ImageURL).
		WithResolution(1920, 1080).
		WithFPS(24).
		WithIncludeCost(true).
		Build()

	fmt.Println("This will take several minutes...")

	response, err := client.VideoInference(ctx, request)
	if err != nil {
		log.Fatalf("Failed to submit video request: %v", err)
	}

	fmt.Printf("\nVideo request submitted: %s\n", response.TaskUUID)
	fmt.Println("Polling for result...")

	finalResp, err := client.PollVideoResult(ctx, response.TaskUUID, 60, 10*time.Second)
	if err != nil {
		log.Fatalf("Failed to get video result: %v", err)
	}

	fmt.Printf("\nVideo generated successfully!\n")
	fmt.Printf("Video UUID: %s\n", finalResp.VideoUUID)
	if finalResp.VideoURL != nil {
		fmt.Printf("Video URL: %s\n", *finalResp.VideoURL)
	}
	if finalResp.Cost != nil {
		fmt.Printf("Cost: $%.4f\n", *finalResp.Cost)
	}
}
