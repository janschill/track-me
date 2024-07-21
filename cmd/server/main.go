package main

import (
	"log"

	"github.com/janschill/track-me/internal/server"
)

func main() {
	port := "8080"
	log.Default().Println("Server starting on port " + port)
	server := server.HttpServer(port)
	log.Fatal(server.ListenAndServe())
}
