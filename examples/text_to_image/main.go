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

	// Create a new client (API key from RUNWARE_API_KEY environment variable)
	client, err := runware.NewClient(nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Connect to the API
	ctx := context.Background()
	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect()

	fmt.Println("Connected to Runware API")

	// Simple text-to-image generation
	prompt := "a serene mountain landscape with a crystal-clear lake reflecting the sky"
	model := "google:4@1"
	width := 1024
	height := 1024

	fmt.Printf("Generating image with prompt: %s\n", prompt)

	response, err := client.TextToImage(ctx, prompt, model, width, height)
	if err != nil {
		log.Fatalf("Failed to generate image: %v", err)
	}

	fmt.Printf("Image generated successfully!\n")
	fmt.Printf("Image UUID: %s\n", response.ImageUUID)
	if response.ImageURL != nil {
		fmt.Printf("Image URL: %s\n", *response.ImageURL)
	}
	if response.Seed != nil {
		fmt.Printf("Seed: %d\n", *response.Seed)
	}
	if response.Cost != nil {
		fmt.Printf("Cost: $%.4f\n", *response.Cost)
	}
}
