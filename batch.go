package runware

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

// batchResult is a generic structure for batch processing results
type batchResult[T any] struct {
	index int
	resp  T
	err   error
}

// processBatch processes multiple requests in parallel using a generic handler
func processBatch[Req any, Resp any](
	ctx context.Context,
	requests []Req,
	handler func(context.Context, Req) (Resp, error),
) ([]Resp, error) {
	if len(requests) == 0 {
		return nil, ErrInvalidRequest
	}

	results := make(chan batchResult[Resp], len(requests))
	var wg sync.WaitGroup

	// Bound concurrency to avoid unbounded goroutines
	maxParallel := runtime.GOMAXPROCS(0) * 4
	if maxParallel < 8 {
		maxParallel = 8
	}
	if len(requests) < maxParallel {
		maxParallel = len(requests)
	}

	sem := make(chan struct{}, maxParallel)

	for i, req := range requests {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, request Req) {
			defer wg.Done()
			defer func() { <-sem }()
			resp, err := handler(ctx, request)
			results <- batchResult[Resp]{index: idx, resp: resp, err: err}
		}(i, req)
	}

	wg.Wait()
	close(results)

	responses := make([]Resp, len(requests))
	var errs []error

	for r := range results {
		if r.err != nil {
			errs = append(errs, fmt.Errorf("request %d: %w", r.index, r.err))
		} else {
			responses[r.index] = r.resp
		}
	}

	if len(errs) > 0 {
		return responses, fmt.Errorf("batch errors: %v", errs)
	}

	return responses, nil
}
