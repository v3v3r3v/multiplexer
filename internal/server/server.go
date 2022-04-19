package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func server(ctx context.Context, addr string, mux *http.ServeMux) (err error) {
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Start http server
	go func() {
		// ListenAndServe returns ErrServerClosed when Shutdown called
		if err = server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("listen error: %s", err)
		}
	}()

	log.Println("Start http server: ", addr)

	// Block until cancellation of context
	<-ctx.Done()

	log.Println("Stop http server: ", addr)

	// Create context with timeout for shutdown of http server
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()

	// Call shutdown
	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server shutdown failed: %s", err)
	}

	if err == http.ErrServerClosed {
		err = nil
	}

	return err
}

func Run(addr string, mux *http.ServeMux) {
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt)
	signal.Notify(osSignal, syscall.SIGTERM)
	signal.Notify(osSignal, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())

	// Wait for signal from OS
	go func() {
		<-osSignal
		cancel()
	}()

	// Start server
	if err := server(ctx, addr, mux); err != nil {
		log.Printf("server error: %s", err)
	}
}