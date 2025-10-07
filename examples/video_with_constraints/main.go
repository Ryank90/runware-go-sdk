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

	// Upload start and end frame images
	fmt.Println("Uploading start frame...")
	startFrame, err := client.UploadImageFromURL(ctx, "https://example.com/start.jpg")
	if err != nil {
		log.Fatalf("Failed to upload start frame: %v", err)
	}

	fmt.Println("Uploading end frame...")
	endFrame, err := client.UploadImageFromURL(ctx, "https://example.com/end.jpg")
	if err != nil {
		log.Fatalf("Failed to upload end frame: %v", err)
	}

	// Create video with frame constraints for smooth transition
	request := runware.NewVideoRequestBuilder(
		"Smooth transition between scenes, cinematic movement",
		"klingai:5@3",
	).
		WithDuration(5).
		WithFirstFrame(startFrame.ImageUUID).
		WithLastFrame(endFrame.ImageUUID).
		WithResolution(1920, 1080).
		WithFPS(24).
		WithIncludeCost(true).
		Build()

	fmt.Println("\nGenerating video with frame constraints...")
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
