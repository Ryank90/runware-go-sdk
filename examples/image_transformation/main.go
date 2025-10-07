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

	// First, generate a base image to transform
	fmt.Println("Generating base image...")
	basePrompt := "a simple portrait of a cat, photograph, neutral background"
	baseResp, err := client.TextToImage(ctx, basePrompt, "runware:101@1", 512, 512)
	if err != nil {
		log.Fatalf("Failed to generate base image: %v", err)
	}

	fmt.Printf("Base image generated with UUID: %s\n", baseResp.ImageUUID)

	// Get image URL for transformation (seedImage requires URL, not just UUID)
	if baseResp.ImageURL == nil {
		log.Fatalf("No image URL returned")
	}

	// Transform the image with a prompt
	prompt := "transform into a watercolor painting, artistic style"
	model := "runware:101@1"
	strength := 0.7

	fmt.Println("Transforming image...")
	response, err := client.ImageToImage(ctx, prompt, model, *baseResp.ImageURL, 512, 512, strength)
	if err != nil {
		log.Fatalf("Failed to transform image: %v", err)
	}

	fmt.Printf("Image transformed successfully!\n")
	fmt.Printf("Result UUID: %s\n", response.ImageUUID)
	if response.ImageURL != nil {
		fmt.Printf("Result URL: %s\n", *response.ImageURL)
	}
}
