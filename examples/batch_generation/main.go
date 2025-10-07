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

	// Create multiple requests
	prompts := []string{
		"a majestic lion in the savanna at sunset",
		"a futuristic cyberpunk cityscape at night",
		"a peaceful zen garden with cherry blossoms",
		"a steampunk airship flying through clouds",
	}

	requests := make([]*runware.ImageInferenceRequest, len(prompts))
	for i, prompt := range prompts {
		requests[i] = runware.NewImageInferenceRequest(prompt, "runware:101@1", 1024, 1024)
		includeCost := true
		requests[i].IncludeCost = &includeCost
	}

	fmt.Printf("Generating %d images in parallel...\n", len(requests))

	// Execute batch generation
	responses, err := client.ImageInferenceBatch(ctx, requests)
	if err != nil {
		log.Fatalf("Batch generation failed: %v", err)
	}

	// Display results
	totalCost := 0.0
	for i, resp := range responses {
		if resp != nil {
			fmt.Printf("\nImage %d: %s\n", i+1, prompts[i])
			fmt.Printf("  UUID: %s\n", resp.ImageUUID)
			if resp.ImageURL != nil {
				fmt.Printf("  URL: %s\n", *resp.ImageURL)
			}
			if resp.Cost != nil {
				fmt.Printf("  Cost: $%.4f\n", *resp.Cost)
				totalCost += *resp.Cost
			}
		}
	}

	fmt.Printf("\nTotal cost: $%.4f\n", totalCost)
}
