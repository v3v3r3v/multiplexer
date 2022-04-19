package main

import (
	"multiplexer/internal/handler"
	"multiplexer/internal/server"
	"net/http"
	"os"
)

func main() {
	addr := os.Getenv("STUB_SERVER_HOST") + ":" + os.Getenv("STUB_SERVER_PORT")

	// Setup routing
	mux := http.NewServeMux()

	mux.Handle("/stub", handler.HttpMethods(
		handler.HandleStub,
		http.MethodGet, http.MethodPost,
	))

	server.Run(addr, mux)
}
