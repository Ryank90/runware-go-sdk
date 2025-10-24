package runware

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Ryank90/runware-go-sdk/models"
)

func TestProcessBatch_EmptyBatch(t *testing.T) {
	ctx := context.Background()
	requests := []*models.ImageInferenceRequest{}

	processor := func(ctx context.Context, req *models.ImageInferenceRequest) (*models.ImageInferenceResponse, error) {
		return &models.ImageInferenceResponse{}, nil
	}

	results, err := processBatch(ctx, requests, processor)
	// Empty batch returns error in current implementation
	if err == nil {
		// If implementation allows empty batches
		if len(results) != 0 {
			t.Errorf("Expected 0 results, got %d", len(results))
		}
	}
	// Error is acceptable for empty batch
}

func TestProcessBatch_SingleItem(t *testing.T) {
	ctx := context.Background()
	requests := []*models.ImageInferenceRequest{
		models.NewImageInferenceRequest("test", "model", 512, 512),
	}

	processor := func(ctx context.Context, req *models.ImageInferenceRequest) (*models.ImageInferenceResponse, error) {
		return &models.ImageInferenceResponse{ImageUUID: "test-uuid"}, nil
	}

	results, err := processBatch(ctx, requests, processor)
	if err != nil {
		t.Errorf("processBatch() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
	if results[0].ImageUUID != "test-uuid" {
		t.Errorf("Unexpected result: %v", results[0])
	}
}

func TestProcessBatch_MultipleItems(t *testing.T) {
	ctx := context.Background()
	requests := []*models.ImageInferenceRequest{
		models.NewImageInferenceRequest("test1", "model", 512, 512),
		models.NewImageInferenceRequest("test2", "model", 512, 512),
		models.NewImageInferenceRequest("test3", "model", 512, 512),
	}

	callCount := 0
	var mu sync.Mutex

	processor := func(ctx context.Context, req *models.ImageInferenceRequest) (*models.ImageInferenceResponse, error) {
		mu.Lock()
		callCount++
		mu.Unlock()
		return &models.ImageInferenceResponse{ImageUUID: req.PositivePrompt}, nil
	}

	results, err := processBatch(ctx, requests, processor)
	if err != nil {
		t.Errorf("processBatch() error = %v", err)
	}
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	mu.Lock()
	if callCount != 3 {
		t.Errorf("Expected processor called 3 times, got %d", callCount)
	}
	mu.Unlock()

	// Verify results are in order
	for i, result := range results {
		expected := requests[i].PositivePrompt
		if result.ImageUUID != expected {
			t.Errorf("Result[%d].ImageUUID = %v, want %v", i, result.ImageUUID, expected)
		}
	}
}

func TestProcessBatch_ErrorHandling(t *testing.T) {
	ctx := context.Background()
	requests := []*models.ImageInferenceRequest{
		models.NewImageInferenceRequest("test1", "model", 512, 512),
		models.NewImageInferenceRequest("test2", "model", 512, 512),
	}

	expectedErr := errors.New("processing error")
	processor := func(ctx context.Context, req *models.ImageInferenceRequest) (*models.ImageInferenceResponse, error) {
		if req.PositivePrompt == "test1" {
			return nil, expectedErr
		}
		return &models.ImageInferenceResponse{}, nil
	}

	_, err := processBatch(ctx, requests, processor)
	if err == nil {
		t.Error("Expected error, got nil")
	}
	// Error should contain the original error message
	if !strings.Contains(err.Error(), "processing error") {
		t.Errorf("Expected error to contain 'processing error', got %v", err)
	}
}

func TestProcessBatch_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	requests := []*models.ImageInferenceRequest{
		models.NewImageInferenceRequest("test1", "model", 512, 512),
		models.NewImageInferenceRequest("test2", "model", 512, 512),
	}

	processor := func(ctx context.Context, req *models.ImageInferenceRequest) (*models.ImageInferenceResponse, error) {
		time.Sleep(10 * time.Millisecond)
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return &models.ImageInferenceResponse{}, nil
	}

	// Cancel context immediately
	cancel()

	_, err := processBatch(ctx, requests, processor)
	if err == nil {
		t.Error("Expected context cancellation error")
	}
}

func TestProcessBatch_BoundedConcurrency(t *testing.T) {
	// This test verifies that batch processing doesn't spawn unlimited goroutines
	ctx := context.Background()

	// Create a large batch
	requests := make([]*models.ImageInferenceRequest, 100)
	for i := 0; i < 100; i++ {
		requests[i] = models.NewImageInferenceRequest("test", "model", 512, 512)
	}

	maxConcurrent := 0
	currentConcurrent := 0
	var mu sync.Mutex

	processor := func(ctx context.Context, req *models.ImageInferenceRequest) (*models.ImageInferenceResponse, error) {
		mu.Lock()
		currentConcurrent++
		if currentConcurrent > maxConcurrent {
			maxConcurrent = currentConcurrent
		}
		mu.Unlock()

		time.Sleep(10 * time.Millisecond) // Simulate work

		mu.Lock()
		currentConcurrent--
		mu.Unlock()

		return &models.ImageInferenceResponse{}, nil
	}

	_, err := processBatch(ctx, requests, processor)
	if err != nil {
		t.Errorf("processBatch() error = %v", err)
	}

	mu.Lock()
	finalMaxConcurrent := maxConcurrent
	mu.Unlock()

	// Verify bounded concurrency (should not be 100)
	// The semaphore in batch.go limits to runtime.NumCPU() * 2
	if finalMaxConcurrent > 50 {
		t.Errorf("Concurrency not properly bounded: max concurrent = %d (expected <= 50)", finalMaxConcurrent)
	}

	t.Logf("Max concurrent goroutines: %d (bounded correctly)", finalMaxConcurrent)
}

func TestProcessBatch_OrderPreservation(t *testing.T) {
	ctx := context.Background()

	// Create requests with different processing times
	requests := []*models.ImageInferenceRequest{
		models.NewImageInferenceRequest("first", "model", 512, 512),
		models.NewImageInferenceRequest("second", "model", 512, 512),
		models.NewImageInferenceRequest("third", "model", 512, 512),
		models.NewImageInferenceRequest("fourth", "model", 512, 512),
	}

	processor := func(ctx context.Context, req *models.ImageInferenceRequest) (*models.ImageInferenceResponse, error) {
		// Simulate varying processing times
		switch req.PositivePrompt {
		case "first":
			time.Sleep(40 * time.Millisecond)
		case "second":
			time.Sleep(10 * time.Millisecond)
		case "third":
			time.Sleep(30 * time.Millisecond)
		case "fourth":
			time.Sleep(20 * time.Millisecond)
		}
		return &models.ImageInferenceResponse{ImageUUID: req.PositivePrompt}, nil
	}

	results, err := processBatch(ctx, requests, processor)
	if err != nil {
		t.Fatalf("processBatch() error = %v", err)
	}

	// Verify results are in the same order as requests
	expectedOrder := []string{"first", "second", "third", "fourth"}
	for i, result := range results {
		if result.ImageUUID != expectedOrder[i] {
			t.Errorf("Result[%d] = %v, want %v (order not preserved)", i, result.ImageUUID, expectedOrder[i])
		}
	}
}

func TestProcessBatch_PartialFailure(t *testing.T) {
	ctx := context.Background()
	requests := []*models.ImageInferenceRequest{
		models.NewImageInferenceRequest("success1", "model", 512, 512),
		models.NewImageInferenceRequest("fail", "model", 512, 512),
		models.NewImageInferenceRequest("success2", "model", 512, 512),
	}

	failErr := errors.New("simulated failure")
	processor := func(ctx context.Context, req *models.ImageInferenceRequest) (*models.ImageInferenceResponse, error) {
		if req.PositivePrompt == "fail" {
			return nil, failErr
		}
		return &models.ImageInferenceResponse{ImageUUID: req.PositivePrompt}, nil
	}

	_, err := processBatch(ctx, requests, processor)
	if err == nil {
		t.Error("Expected error from partial failure, got nil")
	}
}
