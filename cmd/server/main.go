package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/janschill/track-me/internal/server"
)

func main() {
	ping := flag.Bool("ping", false, "Run a self-test and exit")
	flag.Parse()

	if *ping {
		fmt.Println("Self-test passed")
		os.Exit(0)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := server.SetupOTelSDK(ctx)
	if err != nil {
		return
	}
	if err != nil {
			log.Fatalf("Failed to set up OpenTelemetry: %v", err)
	}
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	port := "8080"
	log.Default().Println("Server starting on port " + port)
	srv := server.HttpServer(port, ctx)
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}
	err = srv.Shutdown(context.Background())
}
