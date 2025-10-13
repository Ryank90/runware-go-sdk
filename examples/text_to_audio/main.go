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

	// Text-to-audio generation
	prompt := "Gentle piano melody with soft rain sounds in the background, peaceful and calming"
	model := "elevenlabs:1@1"
	duration := 10 // seconds

	fmt.Printf("\nGenerating audio...\n")
	fmt.Printf("Prompt: %s\n", prompt)
	fmt.Printf("Model: %s\n", model)
	fmt.Printf("Duration: %d seconds\n\n", duration)

	// Build audio request with quality settings
	req := runware.NewAudioRequestBuilder(prompt, model, duration).
		WithAudioSettings(44100, 192). // CD quality, 192kbps
		WithIncludeCost(true).
		Build()

	// Submit the audio generation request (returns quickly with acknowledgment)
	response, err := client.AudioInference(ctx, req)
	if err != nil {
		log.Fatalf("Failed to submit audio request: %v", err)
	}

	fmt.Printf("Audio request submitted successfully!\n")
	fmt.Printf("Task UUID: %s\n\n", response.TaskUUID)
	fmt.Println("Polling for result (this may take 30-60 seconds)...")

	// Poll for the result
	finalResp, err := client.PollAudioResult(ctx, response.TaskUUID, 60, 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to get audio result: %v", err)
	}

	fmt.Printf("\nAudio generated successfully!\n\n")
	if finalResp.AudioUUID != "" {
		fmt.Printf("Audio UUID: %s\n", finalResp.AudioUUID)
	}
	if finalResp.AudioURL != nil {
		fmt.Printf("Audio URL: %s\n", *finalResp.AudioURL)
	}
	if finalResp.Cost != nil {
		fmt.Printf("Cost: $%.4f\n", *finalResp.Cost)
	}
}
