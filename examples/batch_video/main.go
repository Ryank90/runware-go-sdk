package main

import (
	"context"
	"fmt"
	"log"
	"time"

	runware "github.com/Ryank90/runware-go-sdk"
	models "github.com/Ryank90/runware-go-sdk/models"
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

	requests := make([]*models.VideoInferenceRequest, len(prompts))
	for i, prompt := range prompts {
		// Use OpenAI Sora for cost efficiency; requires 1280x720 and durations 4/8/12
		requests[i] = runware.NewVideoRequestBuilder(prompt, "openai:3@1").
			WithDuration(4).
			WithResolution(1280, 720).
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

	// Poll for each result concurrently to minimize wall time
	finalResponses := make([]*models.VideoInferenceResponse, len(responses))
	type idxResp struct {
		idx  int
		resp *models.VideoInferenceResponse
		err  error
	}
	resultCh := make(chan idxResp, len(responses))

	for i, resp := range responses {
		i, resp := i, resp
		if resp == nil {
			continue
		}
		go func() {
			fmt.Printf("\nPolling video %d/%d...\n", i+1, len(responses))
			r, err := client.PollVideoResult(ctx, resp.TaskUUID, 60, 10*time.Second)
			resultCh <- idxResp{idx: i, resp: r, err: err}
		}()
	}

	// Collect all
	for range responses {
		select {
		case r := <-resultCh:
			if r.err != nil {
				log.Printf("Failed to get video %d result: %v", r.idx+1, r.err)
				continue
			}
			finalResponses[r.idx] = r.resp
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
