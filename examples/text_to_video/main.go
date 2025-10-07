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

	// Simple text-to-video generation
	prompt := "A serene beach at sunset with gentle waves crashing on the shore"
	model := "klingai:5@3"
	duration := 5

	fmt.Printf("Generating video: %s\n", prompt)
	fmt.Println("This will take a few minutes...")

	// Submit the video generation request
	response, err := client.TextToVideo(ctx, prompt, model, duration)
	if err != nil {
		log.Fatalf("Failed to submit video request: %v", err)
	}

	fmt.Printf("\nVideo request submitted: %s\n", response.TaskUUID)
	fmt.Println("Polling for result...")

	// Poll for the result
	finalResp, err := client.PollVideoResult(ctx, response.TaskUUID, 60, 10*time.Second)
	if err != nil {
		log.Fatalf("Failed to get video result: %v", err)
	}

	fmt.Printf("\n✓ Video generated successfully!\n")
	fmt.Printf("Video UUID: %s\n", finalResp.VideoUUID)
	if finalResp.VideoURL != nil {
		fmt.Printf("Video URL: %s\n", *finalResp.VideoURL)
	}
	if finalResp.Seed != nil {
		fmt.Printf("Seed: %d\n", *finalResp.Seed)
	}
	if finalResp.Cost != nil {
		fmt.Printf("Cost: $%.4f\n", *finalResp.Cost)
	}
}
