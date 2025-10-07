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

	// Create multiple video generation requests
	prompts := []string{
		"A peaceful forest with sunlight filtering through the trees",
		"Ocean waves crashing on a rocky shore at sunset",
		"A busy city street with cars and people, time-lapse style",
	}

	requests := make([]*runware.VideoInferenceRequest, len(prompts))
	for i, prompt := range prompts {
		requests[i] = runware.NewVideoRequestBuilder(prompt, "klingai:5@3").
			WithDuration(5).
			WithResolution(1920, 1080).
			WithIncludeCost(true).
			Build()
	}

	fmt.Printf("Submitting %d video requests...\n", len(requests))

	// Submit all requests
	responses, err := client.VideoInferenceBatch(ctx, requests)
	if err != nil {
		log.Fatalf("Batch submission failed: %v", err)
	}

	fmt.Println("\nPolling for results...")
	fmt.Println("This will take several minutes...")

	// Poll for each result
	finalResponses := make([]*runware.VideoInferenceResponse, len(responses))
	for i, resp := range responses {
		if resp != nil {
			fmt.Printf("\nPolling video %d/%d...\n", i+1, len(responses))
			finalResp, err := client.PollVideoResult(ctx, resp.TaskUUID, 60, 10*time.Second)
			if err != nil {
				log.Printf("Failed to get video %d result: %v", i+1, err)
				continue
			}
			finalResponses[i] = finalResp
		}
	}

	// Display results
	totalCost := 0.0
	fmt.Println("\n=== Video Generation Results ===")
	for i, resp := range finalResponses {
		if resp != nil {
			fmt.Printf("\nVideo %d: %s\n", i+1, prompts[i])
			fmt.Printf("  UUID: %s\n", resp.VideoUUID)
			if resp.VideoURL != nil {
				fmt.Printf("  URL: %s\n", *resp.VideoURL)
			}
			if resp.Cost != nil {
				fmt.Printf("  Cost: $%.4f\n", *resp.Cost)
				totalCost += *resp.Cost
			}
		}
	}

	fmt.Printf("\nTotal cost: $%.4f\n", totalCost)
}
