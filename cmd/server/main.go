package main

import (
	"context"
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

	port := "8080"
	log.Default().Println("Server starting on port " + port)
	srv := server.HttpServer(port, ctx)
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()
	select {
	case err := <-srvErr:
			log.Fatalf("HTTP server error: %v", err)
	case <-ctx.Done():
			stop()
	}
	if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("HTTP server shutdown error: %v", err)
	}
}
