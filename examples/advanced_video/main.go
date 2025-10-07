package main

import (
	"context"
	"fmt"
	"log"
	"time"

	runware "github.com/Ryank90/runware-go-sdk"
)

func main() {
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

	// Build an advanced video request with provider settings
	request := runware.NewVideoRequestBuilder(
		"A cinematic drone shot flying through a futuristic cyberpunk city at night, neon lights, rain, highly detailed",
		"klingai:5@3",
	).
		WithNegativePrompt("blurry, low quality, static").
		WithDuration(5).
		WithResolution(1920, 1080).
		WithFPS(30).
		WithIncludeCost(true).
		Build()

	fmt.Println("Generating cinematic video with advanced settings...")
	fmt.Println("This will take several minutes...")

	// Submit request
	response, err := client.VideoInference(ctx, request)
	if err != nil {
		log.Fatalf("Failed to submit video request: %v", err)
	}

	fmt.Printf("\nVideo request submitted: %s\n", response.TaskUUID)
	fmt.Println("Polling for result...")

	// Poll for result
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
