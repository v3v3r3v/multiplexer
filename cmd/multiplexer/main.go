package main

import (
	"multiplexer/internal/handler"
	"multiplexer/internal/server"
	"net/http"
	"os"
)

const (
	// Number of incoming requests to be concurrently processed
	concurrentRequestsLimit = 100
)

func main() {
	addr := os.Getenv("MULTIPLEXER_SERVER_HOST") + ":" + os.Getenv("MULTIPLEXER_SERVER_PORT")

	// Setup routing
	mux := http.NewServeMux()
	mux.Handle("/multiplex",
		handler.LimitConcurrentRequests(
		handler.HttpMethods(
			handler.HandleMultiplex,
			http.MethodPost,
		), concurrentRequestsLimit),
	)
	server.Run(addr, mux)
}
