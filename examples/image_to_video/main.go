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

	// Step 1: Generate a base image to animate
	fmt.Println("\nStep 1: Generating base image to animate...")
	imagePrompt := "A serene mountain landscape with a crystal-clear lake, photorealistic, highly detailed"
	imageResp, err := client.TextToImage(ctx, imagePrompt, "runware:101@1", 1024, 576)
	if err != nil {
		log.Fatalf("Failed to generate base image: %v", err)
	}

	fmt.Printf("Base image generated: %s\n", imageResp.ImageUUID)
	if imageResp.ImageURL == nil {
		log.Fatalf("No image URL returned from base image generation")
	}
	fmt.Printf("Image URL: %s\n", *imageResp.ImageURL)

	// Step 2: Generate video from the image
	videoPrompt := "Gentle camera pan across the landscape, realistic clouds moving, water rippling"
	model := "klingai:5@3"
	duration := 5

	fmt.Println("\nStep 2: Generating video from image...")
	fmt.Println("This will take several minutes...")

	// Use the image URL (not UUID) for video generation
	response, err := client.ImageToVideo(ctx, videoPrompt, model, *imageResp.ImageURL, duration)
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
}
