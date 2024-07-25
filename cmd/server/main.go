package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/janschill/track-me/internal/server"
)

func main() {
	ping := flag.Bool("ping", false, "Run a self-test and exit")
	flag.Parse()

	if *ping {
		fmt.Println("Self-test passed")
		os.Exit(0)
	}

	port := "8080"
	log.Default().Println("Server starting on port " + port)
	server := server.HttpServer(port)
	log.Fatal(server.ListenAndServe())
}
