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

	// Upload seed image
	fmt.Println("Uploading seed image...")
	uploadResp, err := client.UploadImageFromURL(ctx, "https://example.com/scene.jpg")
	if err != nil {
		log.Fatalf("Failed to upload image: %v", err)
	}

	fmt.Printf("Image uploaded: %s\n", uploadResp.ImageUUID)

	// Generate video from the image
	prompt := "Animate this scene with gentle camera movement, add realistic motion"
	model := "klingai:5@3"
	duration := 5

	fmt.Println("\nGenerating video from image...")
	fmt.Println("This will take several minutes...")

	response, err := client.ImageToVideo(ctx, prompt, model, uploadResp.ImageUUID, duration)
	if err != nil {
		log.Fatalf("Failed to submit video request: %v", err)
	}

	fmt.Printf("\nVideo request submitted: %s\n", response.TaskUUID)
	fmt.Println("Polling for result...")

	finalResp, err := client.PollVideoResult(ctx, response.TaskUUID, 60, 10*time.Second)
	if err != nil {
		log.Fatalf("Failed to get video result: %v", err)
	}

	fmt.Printf("\n✓ Video generated successfully!\n")
	fmt.Printf("Video UUID: %s\n", finalResp.VideoUUID)
	if finalResp.VideoURL != nil {
		fmt.Printf("Video URL: %s\n", *finalResp.VideoURL)
	}
}
