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

	// Text-to-video generation
	prompt := "A serene beach at sunset with gentle waves crashing on the shore"
	model := "klingai:5@3"
	duration := 5

	fmt.Printf("\nGenerating video...\n")
	fmt.Printf("Prompt: %s\n", prompt)
	fmt.Printf("Model: %s\n", model)
	fmt.Printf("Duration: %d seconds\n\n", duration)

	// Submit the video generation request (returns quickly with acknowledgment)
	response, err := client.TextToVideo(ctx, prompt, model, duration)
	if err != nil {
		log.Fatalf("Failed to submit video request: %v", err)
	}

	fmt.Printf("Video request submitted successfully!\n")
	fmt.Printf("Task UUID: %s\n\n", response.TaskUUID)
	fmt.Println("Polling for result (this may take 2-5 minutes)...")

	// Poll for the result
	finalResp, err := client.PollVideoResult(ctx, response.TaskUUID, 120, 15*time.Second)
	if err != nil {
		log.Fatalf("Failed to get video result: %v", err)
	}

	fmt.Printf("\nVideo generated successfully!\n\n")
	if finalResp.VideoUUID != "" {
		fmt.Printf("Video UUID: %s\n", finalResp.VideoUUID)
	}
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
