package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"shortener/internal/handlers"
	"shortener/internal/redisPkg"
	"sync"
	"time"
)

func AddRoutes(mux *http.ServeMux, clientConn *redisPkg.Database) {
	mux.Handle("/shorten", handlers.HandleShortenURL(clientConn))
	mux.Handle("/", handlers.HandleRedirection(clientConn))
}

func NewServer(redisClient *redisPkg.Database) http.Handler {
	fmt.Println("in new server")
	mux := http.NewServeMux()

	AddRoutes(mux, redisClient)

	// Add middlewares (if any) below
	var handler http.Handler = mux
	// handler = attachLogger(handler)

	return handler
}

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	clientConn := redisPkg.GetClient()
	server := NewServer(clientConn)

	httpServer := &http.Server{
		Addr:    ":8000",
		Handler: server,
	}

	// Start the server in a separate goroutine.
	go func() {
		log.Println("Server is starting...")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		fmt.Println("received shutdown signal")
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)

		defer cancel()
		// want to close redis before the server, so not deferring
		clientConn.Client.Shutdown(context.Background())

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}

		fmt.Println("exited.")
	}()

	wg.Wait()
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
