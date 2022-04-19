package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type FetchResult interface{}

const (
	// Number of requests for fetch data concurrently
	concurrentRequestsLimit = 4
	// Request timeout
	fetchTimeout = 100 * time.Second
)

func FetchUrlList(ctx context.Context, urls []string) ([]FetchResult, error) {

	results := make([]FetchResult, 0, len(urls))

	for res := range makeFetchingQueue(ctx, urls) {
		results = append(results, res)
	}

	if len(results) < len(urls) {
		return nil, fmt.Errorf("an error occurred during the request")
	}

	return results, nil
}

/// Creates queue with limited number of concurrent requests
/// Cancel all requests if any of requests fails
func makeFetchingQueue(ctx context.Context, urls []string) <-chan FetchResult {
	output := make(chan FetchResult)

	go func(ctx context.Context, urls []string, output chan<- FetchResult) {

		// buffered channel for limit of concurrent requests
		concurrentRequests := make(chan struct{}, concurrentRequestsLimit)
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		wg := &sync.WaitGroup{}

		loop:
			for _, url := range urls {
				select {
				case <-ctx.Done():
					break loop
				case concurrentRequests <- struct{}{}:
					wg.Add(1)
					go func(
						ctx context.Context,
						cancel context.CancelFunc,
						wg *sync.WaitGroup,
						concurrentRequests chan struct{},
						output chan<- FetchResult,
						url string,
					) {
						defer wg.Done()
						defer func() { <-concurrentRequests }()

						res, err := fetchUrl(ctx, url)
						if err != nil {
							log.Printf("Fetching error: %v", err)
							// Cancel context to stop all current requests
							cancel()
							return
						}
						output <- res
					}(ctx, cancel, wg, concurrentRequests, output, url)
				}
			}

		// Wait for all tasks complete and close channels
		wg.Wait()
		close(concurrentRequests)
		close(output)
	}(ctx, urls, output)

	return output
}


/// Fetches JSON from single url
/// Request has timeout 1 sec
func fetchUrl(ctx context.Context, url string) (FetchResult, error) {
	httpClient := http.Client{
		Timeout: fetchTimeout,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	var result FetchResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
