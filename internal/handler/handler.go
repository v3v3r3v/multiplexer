package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"multiplexer/internal/fetcher"
	"net/http"
	"strconv"
	"time"
)

const (
	// Default limit is 1000 milliseconds
	defaultMaxSleepMillisecond = 1000
	// Limit number of urls in each request
	maxUrlsCountPerRequest = 20
)

// Test handler with random response time
// Maximum time can be provided in 'limit' query parameter

func HandleStub(w http.ResponseWriter, r *http.Request) {
	maxSleepMillisecond, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		log.Println("defaultMaxSleepMillisecond")
		maxSleepMillisecond = defaultMaxSleepMillisecond
	}
	n := rand.Intn(maxSleepMillisecond)
	time.Sleep(time.Duration(n) * time.Millisecond)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{\"random_number\": %d}", n)))
}

func HandleMultiplex(w http.ResponseWriter, r *http.Request) {
	var urls []string
	err := json.NewDecoder(r.Body).Decode(&urls)

	if err != nil {
		log.Printf("Unable to parse request. Body: %v, err: %v", r.Body, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(urls) > maxUrlsCountPerRequest {
		log.Printf("Unable to process request. The number of urls in the request is more than 20")
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]string{
			"message": fmt.Sprintf("Error: Maximum urls to process %d, given %d", maxUrlsCountPerRequest, len(urls)),
		})
		w.Write(resp)
		return
	}

	ctxFetch, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-r.Context().Done()
		log.Println("Client disconnected")
		cancel()
	}()

	results, err := fetcher.FetchUrlList(ctxFetch, urls)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp, _ := json.Marshal(map[string]string{
			"message": fmt.Sprintf("Error: %v", err),
		})
		w.Write(resp)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp, _ := json.Marshal(results)
	w.Write(resp)
}